package main

import (
	"github.com/go-redis/redis/v7"
	"github.com/kramerdust/weather-reporter/internal/app"
	"github.com/kramerdust/weather-reporter/internal/cache_keeper"
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"log"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("FORECASTER_PORT")
	redisHosts := strings.Split(os.Getenv("REDIS_ADDRS"), ",")

	apiKey, err := forecaster.GetAPIKeyFromEnv()
	if err != nil {
		log.Fatalf("Getting apiKey err: %s", err.Error())
	}
	apiAddr, err := forecaster.GetAPIAddressFromEnv()
	if err != nil {
		log.Fatalf("Getting apiAddr err: %s", err.Error())
	}

	opts := redis.ClusterOptions{
		Addrs: redisHosts,
	}

	ck := cache_keeper.NewRedisCacheKeeper(&opts)

	f := forecaster.WithCacheKeeper(forecaster.NewOWM(apiKey, apiAddr), ck)
	application := app.NewApp(f)
	log.Fatal(application.Run(port))
}
