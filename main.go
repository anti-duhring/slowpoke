package main

import (
	"github.com/anti-duhring/slowpoke/internal/transport-inbound/http"
	"github.com/anti-duhring/slowpoke/internal/transport-outbound/redis"
)

func main() {
	redis, err := redis.Init()
	if err != nil {
		panic(err)
	}

	http.Init(redis)
}
