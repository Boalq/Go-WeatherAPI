package functions

import (
	"WeatherAPI/models"
	"context"
	"log"
	"time"
)

func RedisSetCache(m models.ResponseWeather) {
	ctx := context.Background()

	err := models.Rdb.HSet(ctx, m.City, m).Err()
	if err != nil {
		log.Fatalf("Error")
	}

	ttl := 5

	models.Rdb.Expire(ctx, m.City, time.Duration(ttl)*time.Second)
}
