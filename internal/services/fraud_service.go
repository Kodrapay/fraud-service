package services

import (
	"context"

	"github.com/kodra-pay/fraud-service/internal/dto"
	"github.com/kodra-pay/fraud-service/internal/fraud"
)

// FraudService encapsulates the business logic for fraud detection.
type FraudService struct {
	detector fraud.FraudDetector
}

// NewFraudService creates a new instance of FraudService.
func NewFraudService(detector fraud.FraudDetector) *FraudService {
	return &FraudService{detector: detector}
}

// CheckTransaction orchestrates the fraud check for a given transaction.
func (s *FraudService) CheckTransaction(ctx context.Context, transactionData map[string]interface{}) (fraud.FraudDecision, error) {
	decision, err := s.detector.CheckTransaction(ctx, transactionData)
	if err != nil {
		return fraud.FraudDecision{}, err
	}
	return decision, nil
}

// TrackPaymentLink orchestrates the tracking and analysis of payment links.
func (s *FraudService) TrackPaymentLink(ctx context.Context, linkData map[string]interface{}) (bool, string, error) {
	isSuspicious, reason, err := s.detector.TrackPaymentLink(ctx, linkData)
	if err != nil {
		return false, "", err
	}
	return isSuspicious, reason, nil
}

// ValidatePaymentChannel orchestrates the validation of transactions via payment channels.
func (s *FraudService) ValidatePaymentChannel(ctx context.Context, channelData map[string]interface{}) (bool, string, error) {
	isValid, reason, err := s.detector.ValidatePaymentChannel(ctx, channelData)
	if err != nil {
		return false, "", err
	}
	return isValid, reason, nil
}

func (s *FraudService) Evaluate(_ context.Context, req dto.RiskEvaluateRequest) dto.RiskEvaluateResponse {
	return dto.RiskEvaluateResponse{
		Score:    0.1,
		Decision: "approve",
	}
}

func (s *FraudService) Rules(_ context.Context) []string {
	return []string{"velocity_check", "ip_blacklist"}
}
