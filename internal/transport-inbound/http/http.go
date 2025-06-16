package http

import "github.com/gofiber/fiber/v2"

func Init() *fiber.App {
	app := fiber.New()
	middleware(app)
	routes(app)

	app.Listen(":8080")

	return app
}
