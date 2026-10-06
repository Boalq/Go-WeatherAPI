package models

import (
	"fmt"
)

type ResponseWeather struct {
	City        string  `redis:"city"`
	Country     string  `redis:"country"`
	Temp        float64 `redis:"temp"`
	Description string  `redis:"desc"`
	Feels_like  float64 `redis:"feels_like"`
	Humidity    int     `redis:"humidity"`
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
