package tests

import (
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

var (
	apiKey     = "your API key"
	apiAddress = "http://api.openweathermap.org"
)

func Test_OWM_Forecast(t *testing.T) {
	prepare(t)

	owm := forecaster.NewOWM(apiKey, apiAddress)
	w, err := owm.GetForecast("Moscow")
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err.Error())
	}
	log.Println(w)
}

func Test_OWM_CurWeather(t *testing.T) {
	prepare(t)

	owm:= forecaster.NewOWM(apiKey, apiAddress)
	w, err := owm.GetCurrentWeather("Moscow")
	if err != nil {
		t.Errorf("Unexpected error: %s\n", err.Error())
	}
	log.Println(w)
}

func Test_OWM_BadCityCurWeather(t *testing.T) {
	prepare(t)

	owm := forecaster.NewOWM(apiKey, apiAddress)
	_, err := owm.GetCurrentWeather("Moscowww")
	if !forecaster.IsNotFound(err) {
		t.Errorf("expected not found but got: %s", err.Error())
	}
}

func prepare(t *testing.T) {
	var err error
	apiKey, err = forecaster.GetAPIKeyFromEnv()
	require.NoError(t, err, "apiKey env")

	apiAddress, err = forecaster.GetAPIAddressFromEnv()
	require.NoError(t, err, "apiAddr env")
}