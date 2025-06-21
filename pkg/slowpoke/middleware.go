package slowpoke

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func AddLeakyBucketMiddleware(headerKey string, app *fiber.App, threshold, leakyRate int64, r *redis.Client) {

	app.Use(func(c *fiber.Ctx) error {
		key := c.Get(headerKey)

		b, err := r.Get(c.Context(), key).Result()
		if err == redis.Nil {
			err := r.Set(c.Context(), key, NewLeakyBucket(threshold, leakyRate), 0).Err()
			if err != nil {
				fmt.Errorf("error calling r.Set", err)
				c.Status(500).JSON(fiber.Map{
					"error": "internal server error",
				})
			}
		}
		if err != nil {
			fmt.Errorf("error calling r.Get", err)
			c.Status(500).JSON(fiber.Map{
				"error": "internal server error",
			})
		}

		if !b.CanLeak() {
			return c.Status(429).JSON(fiber.Map{
				"error": "too many requests",
			})
		}

		return c.Next()
	})
}
