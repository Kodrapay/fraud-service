package services

import (
	"context"

	"github.com/kodra-pay/fraud-service/internal/dto"
)

type FraudService struct{}

func NewFraudService() *FraudService { return &FraudService{} }

func (s *FraudService) Evaluate(_ context.Context, req dto.RiskEvaluateRequest) dto.RiskEvaluateResponse {
	return dto.RiskEvaluateResponse{
		Score:    0.1,
		Decision: "approve",
	}
}

func (s *FraudService) Rules(_ context.Context) []string {
	return []string{"velocity_check", "ip_blacklist"}
}
