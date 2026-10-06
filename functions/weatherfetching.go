package functions

import (
	"WeatherAPI/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"

	"golang.org/x/sync/singleflight"
)

func Gettingthatweather(w http.ResponseWriter, r *http.Request) {

	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		log.Fatalf("OPENWEATHER_API .env is not set")
	}

	city := "Berlin"

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	body_reader := strings.NewReader("")
	req, err := http.NewRequest("GET", url, body_reader)

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%s", err)
		return
	}

	cached, err := RedisCaching(fmt.Sprintf("weather:%s", city))
	fmt.Printf("\n----Cached---\n")

	if err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "%s", cached)
		return
	}

	fmt.Printf("\n---NO CACHED---\n")

	var sfGroup singleflight.Group

	data, err, _ := sfGroup.Do(fmt.Sprintf("weather:%s", city), func() (any, error) {

		resp, err := models.HTTPClient.Do(req)
		if err != nil {
			return nil, errors.New("No Respond 404")
		}

		if resp.StatusCode == 400 {
			return "", errors.New("Status Code 400")
		}

		defer resp.Body.Close()

		raw := models.RawParsing{}

		err = json.NewDecoder(resp.Body).Decode(&raw)
		if err != nil {
			return "", errors.New("Json Decoding was unsucessful")
		}

		if len(raw.Weather) < 1 {
			return "", errors.New("Couldn't fetch data")
		}

		better := models.ResponseWeather{
			City:        raw.City,
			Country:     raw.Sys.Country,
			Temp:        math.Round(raw.Main.Temp - 273.15),
			Description: raw.Weather[0].Description,
			Feels_like:  math.Round(raw.Main.Feels_like - 273.15),
			Humidity:    raw.Main.Humidity,
		}

		return better, nil
	})

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%e", err)
		return
	}

	finalResponse, ok := data.(models.ResponseWeather)
	if ok == false {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%e TYPE ASSERTION FAILED", err)
		return
	}

	RedisSetCache(finalResponse)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%s", finalResponse)
}
