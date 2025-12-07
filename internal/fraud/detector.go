package fraud

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/kodra-pay/fraud-service/internal/repository"
)

// FraudDetector defines the interface for fraud detection logic.
type FraudDetector interface {
	CheckTransaction(ctx context.Context, transactionData map[string]interface{}) (FraudDecision, error)
	TrackPaymentLink(ctx context.Context, linkData map[string]interface{}) (bool, string, error)
	ValidatePaymentChannel(ctx context.Context, channelData map[string]interface{}) (bool, string, error)
}

// RuleBasedFraudDetector implements FraudDetector using a set of predefined rules.
type RuleBasedFraudDetector struct {
	repo  repository.FraudDataRepository
	rules []FraudRule
}

// NewRuleBasedFraudDetector creates a new instance of RuleBasedFraudDetector.
func NewRuleBasedFraudDetector(repo repository.FraudDataRepository, rules []FraudRule) *RuleBasedFraudDetector {
	if len(rules) == 0 {
		rules = DefaultRules()
	}
	return &RuleBasedFraudDetector{repo: repo, rules: rules}
}

// CheckTransaction performs fraud checks on transaction data.
func (d *RuleBasedFraudDetector) CheckTransaction(ctx context.Context, transactionData map[string]interface{}) (FraudDecision, error) {
	decision := FraudDecision{
		Decision: "approve",
		Reasons:  []string{},
	}
	var totalScore float64

	for _, rule := range d.rules {
		if !rule.Enabled {
			continue
		}

		// Handle specific predicates that need repository access
		var isTriggered bool
		var err error

		switch rule.ID {
		case "HIGH_VELOCITY_CUSTOMER":
			if customerID, ok := transactionData["customer_id"].(string); ok {
				history, repoErr := d.repo.GetTransactionHistory(ctx, customerID, 24*time.Hour)
				if repoErr != nil {
					return FraudDecision{}, fmt.Errorf("failed to get transaction history for rule %s: %w", rule.ID, repoErr)
				}
				isTriggered = float64(len(history)) > rule.Threshold
			}
		case "SUSPICIOUS_IP_ORIGIN":
			if origin, ok := transactionData["origin"].(string); ok {
				ipData, repoErr := d.repo.GetIPData(ctx, origin)
				if repoErr != nil {
					return FraudDecision{}, fmt.Errorf("failed to get IP data for rule %s: %w", rule.ID, repoErr)
				}
				isTriggered = (ipData != nil && ipData["is_vpn"] == true) || origin == "suspicious_ip" // Placeholder logic
			}
		default:
			isTriggered, err = rule.Predicate(transactionData)
			if err != nil {
				return FraudDecision{}, fmt.Errorf("error evaluating rule %s: %w", rule.ID, err)
			}
		}


		if isTriggered {
			totalScore += rule.ScoreImpact
			decision.Reasons = append(decision.Reasons, rule.Description)
			// Apply immediate decision if rule dictates
			if rule.Decision == "deny" {
				decision.Decision = "deny"
				break // Stop evaluating if denied
			} else if rule.Decision == "flag" && decision.Decision != "deny" {
				decision.Decision = "flag"
			}
		}
	}

	decision.OverallScore = totalScore

	// Final decision based on total score if not already denied
	if decision.Decision != "deny" {
		if totalScore >= 100 { // Example high risk threshold
			decision.Decision = "deny"
		} else if totalScore >= 50 { // Example medium risk threshold
			decision.Decision = "flag"
		}
	}


	return decision, nil
}

// TrackPaymentLink tracks and analyzes payment links for suspicious activity.
func (d *RuleBasedFraudDetector) TrackPaymentLink(ctx context.Context, linkData map[string]interface{}) (bool, string, error) {
	link, ok := linkData["url"].(string)
	if !ok || link == "" {
		return false, "", errors.New("payment link URL is missing or invalid")
	}

	// Example rule: Check for blacklisted domains (placeholder)
	if link == "http://malicious-site.com" {
		return true, fmt.Sprintf("Payment link %s is on a blacklisted domain", link), nil
	}

	// Example rule: Check for suspicious patterns in the URL (e.g., common phishing tactics)
	if containsSuspiciousPattern(link) {
		return true, fmt.Sprintf("Payment link %s contains suspicious patterns", link), nil
	}

	return false, "Payment link appears safe", nil
}

// Placeholder for a function to check suspicious patterns
func containsSuspiciousPattern(link string) bool {
	// In a real scenario, this would involve regex matching, database lookups, etc.
	return false // No suspicious pattern by default for now
}

// ValidatePaymentChannel validates transactions made via different payment channels.
func (d *RuleBasedFraudDetector) ValidatePaymentChannel(ctx context.Context, channelData map[string]interface{}) (bool, string, error) {
	channel, ok := channelData["channel_type"].(string)
	if !ok || channel == "" {
		return false, "", errors.New("payment channel type is missing or invalid")
	}

	transactionID, ok := channelData["transaction_id"].(string)
	if !ok || transactionID == "" {
		return false, "", errors.New("transaction ID is missing or invalid")
	}

	// In a real scenario, this would involve calling out to the respective payment
	// channel service or an external payment gateway to verify the transaction status.
	// For now, it's a placeholder.

	switch channel {
	case "credit_card":
		// Simulate validation logic for credit card
		if transactionID == "fraud_cc_txn_123" {
			return true, fmt.Sprintf("Transaction %s via credit card is fraudulent", transactionID), nil
		}
	case "bank_transfer":
		// Simulate validation logic for bank transfer
		if transactionID == "fraud_bank_txn_456" {
			return true, fmt.Sprintf("Transaction %s via bank transfer is fraudulent", transactionID), nil
		}
	// Add more cases for other payment channels
	default:
		return false, fmt.Sprintf("Unsupported or unknown payment channel: %s", channel), nil
	}

	return false, "Payment channel transaction validated", nil
}
