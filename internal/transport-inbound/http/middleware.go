package http

import (
	"github.com/anti-duhring/slowpoke/pkg/slowpoke"
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func middleware(app *fiber.App, redis *redis.Client) {
	app.Use(userID)

	slowpoke.AddLeakyBucketMiddleware("X-User-Id", app, 10, 10, redis)
}

func userID(c *fiber.Ctx) error {
	userID := c.Get("X-User-Id")

	if userID == "" {

		return c.Status(400).JSON(fiber.Map{
			"error": "X-User-Id not present",
		})
	}

	return c.Next()
}
