package main

import (
	"WeatherAPI/functions"
	"WeatherAPI/models"
	"fmt"
	"log"
	"os"

	"net/http"

	"github.com/joho/godotenv"
)

type Stringer interface {
	String() string
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

	defer models.Rdb.Close()

	city := "Berlin"

	handler := functions.Gettingthatweather(city)
	limitedHandler := functions.RateLimiterMiddleWare(handler)

	http.Handle("/weather", limitedHandler)

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	http.ListenAndServe(port, nil)
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}
