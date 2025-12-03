package handlers

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodra-pay/fraud-service/internal/dto"
	"github.com/kodra-pay/fraud-service/internal/services"
)

type FraudHandler struct {
	svc *services.FraudService
}

func NewFraudHandler(svc *services.FraudService) *FraudHandler { return &FraudHandler{svc: svc} }

func (h *FraudHandler) Evaluate(c *fiber.Ctx) error {
	var req dto.RiskEvaluateRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}
	return c.JSON(h.svc.Evaluate(c.Context(), req))
}

func (h *FraudHandler) Rules(c *fiber.Ctx) error {
	return c.JSON(h.svc.Rules(c.Context()))
}
