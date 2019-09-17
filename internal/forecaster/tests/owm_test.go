package tests

import (
	"log"
	"testing"
	"time"
	"weather-reporter/internal/forecaster"
)

var (
	apiKey     = "your API key"
	apiAddress = "http://api.openweathermap.org"
)

func Test_OWM(t *testing.T) {
	owm := forecaster.NewOWM(apiKey, apiAddress)
	w, err := owm.GetForecast("Moscow", time.Now().Add(time.Hour*3).Unix())
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	log.Println(w)
}
