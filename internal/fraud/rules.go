package fraud

// FraudRule defines the structure for a configurable fraud rule.
type FraudRule struct {
	ID          string
	Description string
	Threshold   float64   // Threshold for the rule to trigger
	ScoreImpact float64   // Score added if this rule is triggered
	Decision    string    // "approve", "flag", "deny" if this rule is triggered
	Enabled     bool
	Predicate   func(transactionData map[string]interface{}) (bool, error) // Function to evaluate the rule
}

// DefaultRules provides a set of example fraud rules.
func DefaultRules() []FraudRule {
	return []FraudRule{
		{
			ID:          "HIGH_AMOUNT_TRANSACTION",
			Description: "Flags transactions with amounts exceeding a high threshold.",
			Threshold:   1000.00,
			ScoreImpact: 50.0,
			Decision:    "flag",
			Enabled:     true,
			Predicate: func(transactionData map[string]interface{}) (bool, error) {
				if amount, ok := transactionData["amount"].(float64); ok {
					return amount > 1000.00, nil
				}
				return false, nil // Or return error if amount is missing/invalid
			},
		},
		{
			ID:          "SUSPICIOUS_IP_ORIGIN",
			Description: "Flags transactions originating from suspicious IP addresses.",
			Threshold:   1.0, // This rule is binary, so threshold of 1 means it either matches or not
			ScoreImpact: 70.0,
			Decision:    "flag",
			Enabled:     true,
			Predicate: func(transactionData map[string]interface{}) (bool, error) {
				if origin, ok := transactionData["origin"].(string); ok {
					// In a real scenario, this would involve a lookup in a blacklist or IP intelligence service
					return origin == "suspicious_ip", nil
				}
				return false, nil
			},
		},
		{
			ID:          "HIGH_VELOCITY_CUSTOMER",
			Description: "Flags customers with unusually high transaction velocity.",
			Threshold:   5.0, // More than 5 transactions in a given lookback period
			ScoreImpact: 60.0,
			Decision:    "flag",
			Enabled:     true,
			Predicate: func(transactionData map[string]interface{}) (bool, error) {
				// This predicate would ideally use the FraudDataRepository to check velocity
				// For now, it's a placeholder that would be integrated with the detector's repo
				return false, nil
			},
		},
		// Add more rules as per industry practices
	}
}

// FraudDecision represents the outcome of a fraud evaluation.
type FraudDecision struct {
	OverallScore float64
	Decision     string // "approve", "flag", "deny"
	Reasons      []string
}
