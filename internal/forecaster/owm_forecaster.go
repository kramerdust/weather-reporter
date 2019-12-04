package forecaster

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"
	"time"
)

const forecastMethod = "/data/2.5/forecast"
const curWeatherMethod = "/data/2.5/weather"

func NewOWM(apiKey, apiAddress string) Forecaster {
	return &owmForecaster{
		httpClient:        &http.Client{},
		apiKey:            apiKey,
		apiAddress:        apiAddress,
		currentWeather:    make(map[string]Weather),
		relevantForecast:  make(map[string][]Weather),
		forecastExpiresAt: make(map[string]time.Time),
		weatherLastCheck:  make(map[string]time.Time),
	}
}

type owmForecaster struct {
	httpClient        *http.Client
	apiKey            string
	apiAddress        string
	relevantForecast  map[string][]Weather
	forecastExpiresAt map[string]time.Time
	currentWeather    map[string]Weather
	weatherLastCheck  map[string]time.Time
}

//easyjson:json
type owmWeatherResponse struct {
	Main main `json:"main"`
}

//easyjson:json
type owmForecastResponse struct {
	List []timedWeather `json:"list"`
}

//easyjson:json
type timedWeather struct {
	Timestamp int64 `json:"dt"`
	Main      main  `json:"main"`
}

//easyjson:json
type main struct {
	Temperature float64 `json:"temp"`
}

func (o *owmForecaster) GetCurrentWeather(city string) (w Weather, err error) {
	req, err := http.NewRequest(http.MethodGet, o.apiAddress+curWeatherMethod, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("q", city)
	q.Add("APPID", o.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = &fcError{
			code: resp.StatusCode,
			msg:  fmt.Sprintf("bad status code: %d", resp.StatusCode),
		}
		return
	}

	data, _ := ioutil.ReadAll(resp.Body)
	wResponse := &owmWeatherResponse{}
	err = json.Unmarshal(data, wResponse)
	if err != nil {
		return
	}

	w.Temperature = kelvinToCelsius(wResponse.Main.Temperature)

	w.Unit = "celsius"

	return
}

func (o *owmForecaster) GetForecast(city string) (weathers []*Weather, err error) {
	weathers, err = o.downloadForecast(city)
	if err != nil {
		return nil, err
	}

	return
}

func (o *owmForecaster) downloadForecast(city string) ([]*Weather, error) {
	req, err := http.NewRequest(http.MethodGet, o.apiAddress+forecastMethod, nil)
	if err != nil {
		return nil, err
	}
	q := url.Values{}
	q.Add("q", city)
	q.Add("APPID", o.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		err = &fcError{
			code: resp.StatusCode,
			msg:  fmt.Sprintf("bad status code: %d", resp.StatusCode),
		}
		return nil, err
	}

	data, _ := ioutil.ReadAll(resp.Body)

	wResponse := &owmForecastResponse{}
	err = json.Unmarshal(data, wResponse)
	if err != nil {
		return nil, err
	}
	if wResponse.List == nil {
		err = errors.New("no forecasts")
		return nil, err
	}

	weathers := make([]*Weather, len(wResponse.List))
	for i := range wResponse.List {
		weathers[i] = &Weather{
			Unit:        "celsius",
			Temperature: kelvinToCelsius(wResponse.List[i].Main.Temperature),
			Timestamp:   wResponse.List[i].Timestamp,
		}
	}

	return weathers, nil
}

func kelvinToCelsius(kelvin float64) int32 {
	return int32(math.Round(kelvin - 273))
}
