# Go Weather API Wrapper for Openweather.com

## prerequisites:
- installed `go`
- installed `redis` (woah caching!!!)

### Usage:
- clone the `repo`
- use your openweathermap.org API (free tier) in `.env` (there is a .env.example template for reference)
- Use `go build .` (currently the City is hard coded in `functions/weatherfetching.go` so change it to your desired City)
- run `redis-server`
- run `./binary` (from the go build)
- curl -s curl "http://localhost:YOURPORT/weather"

Voila! You have a giga hardcoded Weather API Wrapper... (xd)

Project Scope & Idea:
https://roadmap.sh/projects/weather-api-wrapper-service
