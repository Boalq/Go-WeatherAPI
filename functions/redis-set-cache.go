package functions

import (
	"WeatherAPI/models"
	"context"
	"fmt"
	"log"
	"time"
)

func RedisSetCache(m models.ResponseWeather) {
	ctx := context.Background()

	err := models.Rdb.HSet(ctx, fmt.Sprintf("weather:%s", m.City), m).Err()
	if err != nil {
		log.Fatalf("Error")
	}

	ttl := 5

	models.Rdb.Expire(ctx, fmt.Sprintf("weather:%s", m.City), time.Duration(ttl)*time.Second)
}
