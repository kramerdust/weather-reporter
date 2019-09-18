package app

//easyjson:json
type WeatherAPIModel struct {
	City string `json:"city"`
	Unit string `json:"unit"`
	Temperature int `json:"temperature"`
}
