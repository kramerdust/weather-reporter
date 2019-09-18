package forecaster

import (
	"fmt"
	"os"
)

const apiKeyEnv = "FORECASTER_API_KEY"
const apiAddrEnv = "FORECASTER_API_ADDR"

type Weather struct {
	Unit        string
	Temperature int
	Timestamp   int64
}

type Forecaster interface {
	GetCurrentWeather(city string) (Weather, error)
	GetForecast(city string, dt int64) (Weather, error)
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