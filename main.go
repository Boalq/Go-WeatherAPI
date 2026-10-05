package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

var client = &http.Client{
	Timeout: time.Second * 10,
}

type Stringer interface {
	String() string
}

type RawParsing struct {
	City string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
	} `json:"sys"`
	Main struct {
		Temp       float64 `json:"temp"`
		Feels_like float64 `json:"feels_like"`
		Humidity   int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
}

type ResponseWeather struct {
	City        string
	Country     string
	Temp        float64
	Description string
	Feels_like  float64
	Humidity    int
}

func (r ResponseWeather) String() string {
	return fmt.Sprintf("Location: %s | %s \nTemperature: %.1f C \nDescription: %q \nFeels Like: %.1f \nHumilidity: %d",
		r.City,
		r.Country,
		r.Temp,
		r.Description,
		r.Feels_like,
		r.Humidity,
	)
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Couldn't load the .env File!")
	}
}

func main() {
	// What to learn v2:

	// How to implement Redis here.
	// More Deep Dive in http package
	// Understand more what http and http request are


	http.HandleFunc("/weather", gettingthatweather)

	http.ListenAndServe(":8080", nil)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func gettingthatweather(w http.ResponseWriter, req *http.Request) {

	apiKey := os.Getenv("OPENWEATHER_API_KEY")

	if apiKey == "" {
		log.Fatalf("OPENWEATHER_API .env is not set")
	}

	city := "Berlin"

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?q=%s&appid=%s", city, apiKey)

	body_reader := strings.NewReader("")
	req, err := http.NewRequest("GET", url, body_reader)
	if err != nil {
		log.Fatal(err)
		return
	}


	resp, err := client.Do(req)
	if err != nil {
		fmt.Fprint(w, "No Response")
		return
	}

	if resp.StatusCode == 400 {
		fmt.Fprintf(w, "%s country is not real or no data to it", city)
		return
	}

	defer resp.Body.Close()

	raw := RawParsing{}

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil{
		fmt.Fprintf(w, "ERROR")
		return
	}

	if len(raw.Weather) < 1{
		fmt.Fprintf(w, "NO DATA")
		return
	}

	better := ResponseWeather{
		City:        raw.City,
		Country:     raw.Sys.Country,
		Temp:        math.Round(raw.Main.Temp - 273.15),
		Description: raw.Weather[0].Description,
		Feels_like:  math.Round(raw.Main.Feels_like - 273.15),
		Humidity:    raw.Main.Humidity,
	}

	fmt.Fprintf(w, "%s", better)
}
