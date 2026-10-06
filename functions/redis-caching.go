package functions

import (
	"WeatherAPI/models"
	"context"
	"errors"
	"strings"
)

func RedisCaching(key string) (models.ResponseWeather, error) {
	ctx := context.Background()

	var cached models.ResponseWeather
	err := models.Rdb.HGetAll(ctx, key).Scan(&cached)
	if err != nil {
		panic(err)
	}

	if strings.Compare(cached.City, "") == 0 {
		return models.ResponseWeather{}, errors.New("Nothing in cache")
	}

	return cached, nil
}
