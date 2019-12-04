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
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

type Application struct {
	fc forecaster.Forecaster
	r  *router.Router

	promStatuses *prometheus.CounterVec
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
	a.promStatuses = prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "forecaster",
		Name:      "status_code_counter",
	},
		[]string{"method", "code"},
	)
	prometheus.MustRegister(a.promStatuses)

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

		a.promStatuses.WithLabelValues(string(ctx.Method()), strconv.Itoa(ctx.Response.StatusCode())).Inc()

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

	weathers, err := a.fc.GetForecast(city)
	if err != nil {
		if forecaster.IsNotFound(err) {
			ctx.SetStatusCode(fasthttp.StatusNotFound)
			return
		}
		log.Println(err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		return
	}

	weather := weathers[forecaster.FindClosest(weathers, int64(dt))]

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
