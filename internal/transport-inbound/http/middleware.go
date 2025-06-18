package http

import (
	"github.com/anti-duhring/slowpoke/pkg/slowpoke"
	"github.com/gofiber/fiber/v2"
)

func middleware(app *fiber.App) {
	app.Use(userID)

	b := map[string]slowpoke.Bucket{}

	app.Use(func(c *fiber.Ctx) error {
		return rateLimit(c, b)
	})
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

func rateLimit(c *fiber.Ctx, buckets map[string]slowpoke.Bucket) error {
	userID := c.Get("X-User-Id")

	if _, ok := buckets[userID]; !ok {
		buckets[userID] = slowpoke.NewBucket(10, 10)
	}

	b := buckets[userID]
	if !b.HasToken() {
		return c.Status(429).JSON(fiber.Map{
			"error": "too many requests",
		})
	}

	return c.Next()
}
