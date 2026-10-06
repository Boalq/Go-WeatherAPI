package functions

import(
	"os"
	"log"
	"net/http"
	"fmt"
	"strings"
	"WeatherAPI/models"
	"math"
	"encoding/json"
)

func Gettingthatweather(city string) http.Handler {

	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		log.Fatalf("OPENWEATHER_API .env is not set")
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	body_reader := strings.NewReader("")
	req, err := http.NewRequest("GET", url, body_reader)
	
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%s", err)
		})
	}

	cached, err := RedisCaching(city)
	
	if err == nil{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusAccepted)
			fmt.Fprintf(w, "%s", cached)
		})
	}


	resp, err := models.HTTPClient.Do(req)
	if err != nil {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%e", err)
		})
	}

	if resp.StatusCode == 400 {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%e", err)
		})
	}

	defer resp.Body.Close()

	raw := models.RawParsing{}

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%e", err)
		})
	}

	if len(raw.Weather) < 1{
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "%e", err)
		})
	}

	better := models.ResponseWeather{
		City:        raw.City,
		Country:     raw.Sys.Country,
		Temp:        math.Round(raw.Main.Temp - 273.15),
		Description: raw.Weather[0].Description,
		Feels_like:  math.Round(raw.Main.Feels_like - 273.15),
		Humidity:    raw.Main.Humidity,
	}

	RedisSetCache(better)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusAccepted)
		fmt.Fprintf(w, "%s", better)
	})
}