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
		relevantForecast:  make(map[string][]Weather),
		forecastExpiresAt: make(map[string]time.Time),
	}
}

type owmForecaster struct {
	httpClient        *http.Client
	apiKey            string
	apiAddress        string
	relevantForecast  map[string][]Weather
	forecastExpiresAt map[string]time.Time
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
		err = fmt.Errorf("bad status code: %d", resp.StatusCode)
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

func (o *owmForecaster) GetForecast(city string, dt int64) (w Weather, err error) {
	var weathers []Weather
	if o.forecastExpiresAt[city].Before(time.Now()) {
		err = o.downloadForecast(city)
		if err != nil {
			return
		}
	}

	weathers = o.relevantForecast[city]
	if weathers == nil {
		err = o.downloadForecast(city)
		if err != nil {
			return
		}
	}

	w.Temperature = weathers[findClosest(o.relevantForecast[city], dt)].Temperature
	w.Unit = "celsius"

	return
}

func (o *owmForecaster) downloadForecast(city string) error {
	req, err := http.NewRequest(http.MethodGet, o.apiAddress+forecastMethod, nil)
	if err != nil {
		return err
	}
	q := url.Values{}
	q.Add("q", city)
	q.Add("APPID", o.apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := o.httpClient.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d", resp.StatusCode)
		return err
	}

	data, _ := ioutil.ReadAll(resp.Body)

	wResponse := &owmForecastResponse{}
	err = json.Unmarshal(data, wResponse)
	if err != nil {
		return err
	}
	if wResponse.List == nil {
		err = errors.New("no forecasts")
		return err
	}

	weathers := make([]Weather, len(wResponse.List))
	for i := range wResponse.List {
		weathers[i] = Weather{
			Unit:        "celsius",
			Temperature: kelvinToCelsius(wResponse.List[i].Main.Temperature),
			Timestamp:   wResponse.List[i].Timestamp,
		}
	}

	o.relevantForecast[city] = weathers
	o.forecastExpiresAt[city] = time.Now().Add(time.Hour * 24 * 5)

	return nil
}

func findClosest(slice []Weather, toFind int64) int {
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

func kelvinToCelsius(kelvin float64) int {
	return int(math.Round(kelvin - 273))
}
