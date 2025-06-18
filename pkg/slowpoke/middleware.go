package slowpoke

import "github.com/gofiber/fiber/v2"

func AddLeakyBucketMiddleware(headerKey string, app *fiber.App, threshold, leakyRate int64) {

	b := map[string]LeakyBucket{}

	app.Use(func(c *fiber.Ctx) error {
		key := c.Get(headerKey)

		if _, ok := b[key]; !ok {
			b[key] = NewLeakyBucket(threshold, leakyRate)
		}

		if !b[key].CanLeak() {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
			})
		}

		return c.Next()
	})
}
