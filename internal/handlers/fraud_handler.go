package handlers

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/kodra-pay/fraud-service/internal/services"
)

// FraudAPIHandler handles HTTP requests related to fraud detection.
type FraudAPIHandler struct {
	svc *services.FraudService
}

// NewFraudAPIHandler creates a new instance of FraudAPIHandler.
func NewFraudAPIHandler(svc *services.FraudService) *FraudAPIHandler {
	return &FraudAPIHandler{svc: svc}
}

// CheckTransaction handles the request to check a transaction for fraud.
func (h *FraudAPIHandler) CheckTransaction(c *fiber.Ctx) error {
	var transactionData map[string]interface{}
	if err := c.BodyParser(&transactionData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Manual input validation
	customerID, ok := transactionData["customer_id"].(string)
	if !ok || customerID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "customer_id is required and must be a string")
	}
	amount, ok := transactionData["amount"].(float64)
	if !ok || amount <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "amount is required and must be a positive number")
	}
	currency, ok := transactionData["currency"].(string)
	if !ok || currency == "" {
		return fiber.NewError(fiber.StatusBadRequest, "currency is required and must be a string")
	}

	decision, err := h.svc.CheckTransaction(c.Context(), transactionData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(decision)
}

// TrackPaymentLink handles the request to track a payment link for suspicious activity.
func (h *FraudAPIHandler) TrackPaymentLink(c *fiber.Ctx) error {
	var linkData struct {
		URL string `json:"url"`
	}
	if err := c.BodyParser(&linkData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	if linkData.URL == "" {
		return fiber.NewError(fiber.StatusBadRequest, "url is required")
	}

	isSuspicious, reason, err := h.svc.TrackPaymentLink(c.Context(), map[string]interface{}{"url": linkData.URL})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"is_suspicious": isSuspicious,
		"reason":        reason,
	})
}

// ValidatePaymentChannel handles the request to validate a transaction via a payment channel.
func (h *FraudAPIHandler) ValidatePaymentChannel(c *fiber.Ctx) error {
	var channelData map[string]interface{}
	if err := c.BodyParser(&channelData); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid request body")
	}

	// Manual input validation
	channelType, ok := channelData["channel_type"].(string)
	if !ok || channelType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "channel_type is required and must be a string")
	}
	transactionID, ok := channelData["transaction_id"].(string)
	if !ok || transactionID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "transaction_id is required and must be a string")
	}

	isValid, reason, err := h.svc.ValidatePaymentChannel(c.Context(), channelData)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(fiber.Map{
		"is_valid": isValid,
		"reason":   reason,
	})
}

// GetTransactionDetails handles the request to get transaction details by reference.
func (h *FraudAPIHandler) GetTransactionDetails(c *fiber.Ctx) error {
	reference := c.Params("reference")
	if reference == "" {
		return fiber.NewError(fiber.StatusBadRequest, "transaction reference is required")
	}

	transaction, err := h.svc.GetTransactionDetailsByReference(c.Context(), reference)
	if err != nil {
		// Differentiate between "not found" and other errors
		if err.Error() == fmt.Sprintf("transaction with reference %s not found", reference) {
			return fiber.NewError(fiber.StatusNotFound, err.Error())
		}
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	return c.JSON(transaction)
}

// Register adds the fraud-related routes to the Fiber app.
func (h *FraudAPIHandler) Register(app *fiber.App) {
	fraudGroup := app.Group("/fraud")
	fraudGroup.Use(limiter.New(limiter.Config{
		Max:        5,               // Allow 5 requests
		Expiration: 1 * time.Second, // per 1 second
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.Get("X-API-Key") // Rate limit per API key
		},
	}))
	fraudGroup.Post("/check-transaction", h.CheckTransaction)
	fraudGroup.Post("/track-payment-link", h.TrackPaymentLink)
	fraudGroup.Post("/validate-payment-channel", h.ValidatePaymentChannel)
	fraudGroup.Get("/transactions/:reference", h.GetTransactionDetails) // New route
}
