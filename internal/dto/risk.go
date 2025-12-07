package dto

type RiskEvaluateRequest struct {
	TransactionReference int     `json:"transaction_reference"`
	Amount               float64 `json:"amount"`
	Currency             string  `json:"currency"`
	CustomerID           int     `json:"customer_id"`
	MerchantID           int     `json:"merchant_id"`
}

type RiskEvaluateResponse struct {
	Score    float64 `json:"score"`
	Decision string  `json:"decision"`
}
