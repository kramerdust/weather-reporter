package main

import (
	"github.com/kramerdust/weather-reporter/internal/app"
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"log"
	"os"
)

func main() {
	port := os.Getenv("FORECASTER_PORT")
	apiKey, err := forecaster.GetAPIKeyFromEnv()
	if err != nil {
		log.Fatalf("Getting apiKey err: %s", err.Error())
	}
	apiAddr, err := forecaster.GetAPIAddressFromEnv()
	if err != nil {
		log.Fatalf("Getting apiAddr err: %s", err.Error())
	}

	f := forecaster.NewOWM(apiKey, apiAddr)
	application := app.NewApp(f)
	application.Init()
	log.Fatal(application.Run(port))
}
