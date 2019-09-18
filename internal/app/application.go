package app

import (
	"encoding/json"
	"github.com/fasthttp/router"
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"github.com/valyala/fasthttp"
	"log"
	"unicode/utf8"
)


type Application struct {
	fc forecaster.Forecaster
	r  *router.Router
}

func NewApp(fc forecaster.Forecaster) *Application {
	r := router.New()
	return &Application{
		fc: fc,
		r:  r,
	}
}

func (a *Application) Init() {
	a.r.GET("/v1/forecast", a.GetForecast)
}

func (a *Application) Run(port string) error {
	return fasthttp.ListenAndServe(":" + port, a.r.Handler)
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