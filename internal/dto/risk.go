package dto

type RiskEvaluateRequest struct {
	TransactionReference string  `json:"transaction_reference"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	CustomerID           string  `json:"customer_id"`
	MerchantID           string  `json:"merchant_id"`
}

type RiskEvaluateResponse struct {
	Score    float64 `json:"score"`
	Decision string  `json:"decision"`
}
