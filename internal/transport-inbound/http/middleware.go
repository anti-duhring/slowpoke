package http

import (
	"github.com/anti-duhring/slowpoke/pkg/slowpoke"
	"github.com/gofiber/fiber/v2"
)

func middleware(app *fiber.App) {
	app.Use(userID)

	slowpoke.AddLeakyBucketMiddleware("X-User-Id", app, 10, 10)
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
