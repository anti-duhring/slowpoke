package http

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func Init(redis *redis.Client) *fiber.App {
	app := fiber.New()
	middleware(app, redis)
	routes(app)

	port := "8080"
	p := os.Getenv("API_PORT")
	if p != "" {
		port = p
	}

	app.Listen(fmt.Sprintf(":%v", port))

	return app
}
