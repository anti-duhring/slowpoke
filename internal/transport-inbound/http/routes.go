package http

import "github.com/gofiber/fiber/v2"

func routes(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		userID := c.Get("X-User-Id")

		return c.JSON(fiber.Map{
			"message": "hello world",
			"userID":  userID,
		})
	})
}
