package middleware

import (
	"github.com/gofiber/fiber/v2"
)

// APIKeyAuth is a middleware that checks for a valid API key.
func APIKeyAuth(apiKey string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Allow requests to /health without API key
		if c.Path() == "/health" {
			return c.Next()
		}

		providedAPIKey := c.Get("X-API-Key")
		if providedAPIKey == "" {
			return fiber.NewError(fiber.StatusUnauthorized, "API key missing")
		}

		if providedAPIKey != apiKey {
			return fiber.NewError(fiber.StatusUnauthorized, "Invalid API key")
		}

		return c.Next()
	}
}