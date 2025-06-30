package slowpoke

import (
	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AddLeakyBucketMiddleware(headerKey string, app *fiber.App, threshold, leakyRate int64, r *redis.Client) {
	b := NewRedisLeakyBucket(threshold, leakyRate, r)

	app.Use(func(c *fiber.Ctx) error {
		key := c.Get(headerKey)

		allow, err := b.Allow(key, c.Context())
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		if !allow {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
			})
		}

		return c.Next()
	})
}
