package models

import(
	"net/http"
	"time"
	"github.com/redis/go-redis/v9"
)


var HTTPClient = &http.Client{
	Timeout: time.Second * 10,
}


var Rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", 
	DB:       0,  
	Protocol: 2,
})