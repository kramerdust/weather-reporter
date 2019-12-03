package app

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"log"
	"strings"
	"time"
	"unicode/utf8"
)

type Application struct {
	fc forecaster.Forecaster
	r  *router.Router

	promOK       prometheus.Counter
	promNotFound prometheus.Counter
	promIntErr   prometheus.Counter
	promRespTime prometheus.Histogram
}

func NewApp(fc forecaster.Forecaster) *Application {
	// log.SetOutput(os.Stdout)
	r := router.New()

	a := &Application{
		fc: fc,
		r:  r,
	}

	a.initHandlers()
	a.initMetrics()
	return a
}

func (a *Application) initHandlers() {
	a.r.GET("/v1/forecast", a.GetForecast)
	a.r.GET("/v1/current", a.GetCurrent)
	a.r.GET("/metrics", fasthttpadaptor.NewFastHTTPHandler(promhttp.Handler()))
}

func (a *Application) initMetrics() {
	a.promOK = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "forecaster_http_200",
	})
	a.promNotFound = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "forecaster_http_404",
	})
	a.promIntErr = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "forecaster_http_500",
	})
	prometheus.MustRegister(a.promNotFound, a.promOK, a.promIntErr)

	a.promRespTime = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "forecaster_response_time",
		Buckets: prometheus.DefBuckets,
	})
	prometheus.MustRegister(a.promRespTime)
}

func (a *Application) metrics(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		if strings.HasSuffix(string(ctx.Path()), "metrics") {
			next(ctx)
			return
		}

		start := time.Now()

		next(ctx)

		switch ctx.Response.StatusCode() {
		case fasthttp.StatusOK:
			a.promOK.Inc()
		case fasthttp.StatusNotFound:
			a.promNotFound.Inc()
		case fasthttp.StatusInternalServerError:
			a.promIntErr.Inc()
		}

		total := time.Now().Sub(start).Seconds()
		a.promRespTime.Observe(total)
	}
}

func (a *Application) Run(port string) error {
	log.Println("Starting service ...")
	return fasthttp.ListenAndServe(":"+port, a.metrics(a.r.Handler))
}

func (a *Application) GetForecast(ctx *fasthttp.RequestCtx) {
	cityBytes := ctx.Request.URI().QueryArgs().Peek("city")
	if !utf8.Valid(cityBytes) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	city := string(cityBytes)

	dt, err := ctx.Request.URI().QueryArgs().GetUint("dt")
	if err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}

	weather, err := a.fc.GetForecast(city, int64(dt))
	if err != nil {
		if forecaster.IsNotFound(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	apiModel := WeatherAPIModel{
		City:        city,
		Unit:        weather.Unit,
		Temperature: weather.Temperature,
	}

	data, err := json.Marshal(apiModel)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}
	log.Printf("Response: %s\nClient: %s\n", string(data), ctx.RemoteIP().String())
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Add(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	ctx.Response.BodyWriter().Write(data)
}

func (a *Application) GetCurrent(ctx *fasthttp.RequestCtx) {
	cityBytes := ctx.Request.URI().QueryArgs().Peek("city")
	if !utf8.Valid(cityBytes) {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		return
	}
	city := string(cityBytes)

	weather, err := a.fc.GetCurrentWeather(city)
	if err != nil {
		if forecaster.IsNotFound(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	apiModel := WeatherAPIModel{
		City:        city,
		Unit:        weather.Unit,
		Temperature: weather.Temperature,
	}

	data, err := json.Marshal(apiModel)
	if err != nil {
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.Response.Header.Add(fasthttp.HeaderContentType, "application/json; charset=utf-8")
	ctx.Response.BodyWriter().Write(data)
}
