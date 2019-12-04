package app

//easyjson:json
type WeatherAPIModel struct {
	City string `json:"city"`
	Unit string `json:"unit"`
	Temperature int32 `json:"temperature"`
}
