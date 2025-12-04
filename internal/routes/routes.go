package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/fraud-service/internal/handlers"
)

func Register(app *fiber.App, serviceName string) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)
}
