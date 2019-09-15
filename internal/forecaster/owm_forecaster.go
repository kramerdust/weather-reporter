package forecaster

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"net/url"
)

const forecastMethod = "/data/2.5/forecast"
const curWeathermethod = "/data/2.5/weather"

func NewOWM(apiKey, apiAddress string) Forecaster {
	return &owmForecaster{
		httpClient: &http.Client{},
		apiKey:     apiKey,
		apiAddress: apiAddress,
	}
}

type owmForecaster struct {
	httpClient *http.Client
	apiKey string
	apiAddress string
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
	Main main `json:"main"`
}

//easyjson:json
type main struct {
	Temperature float64 `json:"temp"`
}

func (o *owmForecaster) GetCurrentWeather(city string) (w Weather, err error) {
	req, err := http.NewRequest(http.MethodGet, o.apiAddress + curWeathermethod, nil)
	if err != nil {
		return
	}
	q := req.URL.Query()
	q.Add("q", city)
	q.Add("APPID", o.apiKey)
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status code: %d", resp.StatusCode)
	}

	data,_ := ioutil.ReadAll(resp.Body)
	wResponse := &owmWeatherResponse{}
	err = json.Unmarshal(data, wResponse)
	if err != nil {
		return
	}

	w.City = city
	w.Temperature = kelvinToCelsius(wResponse.Main.Temperature)

	w.Unit = "celsius"

	return
}



func (o *owmForecaster) GetForecast(city string, dt int64) (w Weather, err error) {
	req, err := http.NewRequest(http.MethodGet, o.apiAddress + forecastMethod, nil)
	if err != nil {
		return
	}
	q := url.Values{}
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

	data,_ := ioutil.ReadAll(resp.Body)
	log.Println(string(data))
	wResponse := &owmForecastResponse{}
	err = json.Unmarshal(data, wResponse)
	if err != nil {
		return
	}
	if wResponse.List == nil {
		err = errors.New("no forecasts")
	}

	w.City = city
	w.Temperature = kelvinToCelsius(wResponse.List[findClosest(wResponse.List, dt)].Main.Temperature)
	w.Unit = "celsius"

	return
}

func findClosest(slice []timedWeather, toFind int64) int {
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

func kelvinToCelsius (kelvin float64) int {
	return int(math.Round(kelvin - 273))
}