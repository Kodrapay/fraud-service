package routes

import (
	"os" // Added import for os

	"github.com/gofiber/fiber/v2"
	"github.com/kodra-pay/fraud-service/internal/fraud"
	"github.com/kodra-pay/fraud-service/internal/handlers"
	"github.com/kodra-pay/fraud-service/internal/repository" // Add repository import
	"github.com/kodra-pay/fraud-service/internal/services"
)

func Register(app *fiber.App, serviceName string) {
	health := handlers.NewHealthHandler(serviceName)
	health.Register(app)

	// Initialize Fraud components
	fraudRepo := repository.NewInMemoryFraudDataRepository()                          // Initialize in-memory repo
	fraudDetector := fraud.NewRuleBasedFraudDetector(fraudRepo, fraud.DefaultRules()) // Pass repo and rules
	
	transactionServiceURL := os.Getenv("TRANSACTION_SERVICE_URL")
	fraudService := services.NewFraudService(fraudDetector, transactionServiceURL)
	
	fraudAPIHandler := handlers.NewFraudAPIHandler(fraudService)
	fraudAPIHandler.Register(app)
}
