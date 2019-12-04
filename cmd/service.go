package main

import (
	"github.com/chasex/redis-go-cluster"
	"github.com/kramerdust/weather-reporter/internal/app"
	"github.com/kramerdust/weather-reporter/internal/cache_keeper"
	"github.com/kramerdust/weather-reporter/internal/forecaster"
	"log"
	"os"
	"strings"
	"time"
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

	opts := redis.Options{
		StartNodes:   redisHosts,
		ConnTimeout: 50 * time.Millisecond,
		ReadTimeout: 50 * time.Millisecond,
		WriteTimeout: 50 * time.Millisecond,
		KeepAlive: 16,
		AliveTime: 60 * time.Second,
	}

	ck, err := cache_keeper.NewRedisCacheKeeper(&opts)
	if err != nil {
		log.Fatalf("Connecting to redis cluster err: %s", err.Error())
	}

	f := forecaster.WithCacheKeeper(forecaster.NewOWM(apiKey, apiAddr), ck)
	application := app.NewApp(f)
	log.Fatal(application.Run(port))
}
