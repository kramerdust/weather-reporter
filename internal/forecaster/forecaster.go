package forecaster

import (
	"fmt"
	"github.com/kramerdust/weather-reporter/internal/cache_keeper"
	"log"
	"os"
)

const apiKeyEnv = "FORECASTER_API_KEY"
const apiAddrEnv = "FORECASTER_API_ADDR"


type Forecaster interface {
	GetCurrentWeather(city string) (Weather, error)
	GetForecast(city string) ([]*Weather, error)
}

func GetAPIKeyFromEnv() (string, error) {
	return getEnv(apiKeyEnv)
}

func GetAPIAddressFromEnv() (string, error) {
	return getEnv(apiAddrEnv)
}

func getEnv(key string) (string, error) {
	if val, found := os.LookupEnv(key); found {
		return val, nil
	} else {
		return "", fmt.Errorf("no env with name: %s", key)
	}
}

func WithCacheKeeper(forecaster Forecaster, ck cache_keeper.CacheKeeper) Forecaster {
	return forecasterWithCache{
		ck:       ck,
		internal: forecaster,
	}
}

type forecasterWithCache struct {
	ck cache_keeper.CacheKeeper
	internal Forecaster
}

func (f forecasterWithCache) GetCurrentWeather(city string) (Weather, error) {
	var w Weather
	err := f.ck.Get(city + "_current", &w)
	if cache_keeper.CacheNotFound(err) {
		return f.getAndCacheCurrent(city)
	}

	return w, nil
}

func (f forecasterWithCache) GetForecast(city string) ([]*Weather, error) {
	forecast := Forecast{Weathers:make([]*Weather,0)}

	err := f.ck.Get(city + "_forecast", &forecast)
	if cache_keeper.CacheNotFound(err) {
		return f.getAndCacheForecast(city)
	}

	return forecast.Weathers, nil

}

func (f forecasterWithCache) getAndCacheCurrent(city string) (Weather, error) {
	w, err := f.internal.GetCurrentWeather(city)
	if err != nil {
		return Weather{}, err
	}

	err = f.ck.Set(city + "_current", &w)
	if err != nil {
		log.Printf("failed to cache weather :%s\n", err.Error())
	}

	return w, nil
}

func (f forecasterWithCache) getAndCacheForecast(city string) ([]*Weather, error) {
	weathers, err := f.internal.GetForecast(city)
	if err != nil {
		return nil, err
	}

	err = f.ck.Set(city + "_current", &Forecast{Weathers:weathers})
	if err != nil {
		log.Printf("failed to cache forecast :%s\n", err.Error())
	}

	return weathers, nil
}

func FindClosest(slice []*Weather, toFind int64) int {
	b := 0
	e := len(slice) - 1
	var m int
	for b < e {
		m = (e + b) / 2
		if toFind == slice[m].Timestamp {
			return m
		}
		if toFind > slice[m].Timestamp {
			b = m + 1
		} else {
			e = m
		}
	}
	return b
}