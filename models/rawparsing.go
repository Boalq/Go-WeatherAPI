package models

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
