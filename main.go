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

var middleware = []func(http.HandlerFunc) http.HandlerFunc{
	functions.RateLimiterMiddleWare,
}

func main() {
	// What to learn v2:

	// More Deep Dive in http package
	// Understand more what http and http request are

	defer models.Rdb.Close()

	handler := functions.Gettingthatweather

	for _, m := range middleware {
		handler = m(handler)
	}

	http.HandleFunc("/weather", handler)

	port := fmt.Sprintf(":%v", os.Getenv("PORT"))
	http.ListenAndServe(port, nil)
}
