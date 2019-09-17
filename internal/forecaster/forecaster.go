package forecaster

type Weather struct {
	Unit        string
	Temperature int
	Timestamp   int64
}

type Forecaster interface {
	GetCurrentWeather(city string) (Weather, error)
	GetForecast(city string, dt int64) (Weather, error)
}
