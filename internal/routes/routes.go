package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/fraud-service/internal/handlers"
	"github.com/kodra-pay/fraud-service/internal/services"
)

func Register(app *fiber.App, service string) {
	health := handlers.NewHealthHandler(service)
	health.Register(app)

	svc := services.NewFraudService()
	h := handlers.NewFraudHandler(svc)
	api := app.Group("/risk")
	api.Post("/evaluate", h.Evaluate)
	api.Get("/rules", h.Rules)
}
