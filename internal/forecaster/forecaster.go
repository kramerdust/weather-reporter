package forecaster

type Weather struct {
	City string
	Unit string
	Temperature int
}

type Forecaster interface {
	GetCurrentWeather(city string) (Weather, error)
	GetForecast(city string, dt int64) (Weather, error)
}

