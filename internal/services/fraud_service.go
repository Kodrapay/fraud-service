package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/kodra-pay/fraud-service/internal/dto"
	"github.com/kodra-pay/fraud-service/internal/fraud"
)

// TransactionResponse DTO for returning transaction information (copied from transaction-service/internal/dto/dto.go)
type TransactionResponse struct {
	ID            int       `json:"id"`
	Reference     string    `json:"reference"`
	MerchantID    int       `json:"merchant_id"`
	CustomerEmail string    `json:"customer_email"`
	CustomerID    int       `json:"customer_id"`
	CustomerName  string    `json:"customer_name,omitempty"`
	Amount        int64     `json:"amount"`
	Currency      string    `json:"currency"`
	Status        string    `json:"status"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
}

// FraudService encapsulates the business logic for fraud detection.

type FraudService struct {
	detector fraud.FraudDetector

	transactionServiceURL string
}

// NewFraudService creates a new instance of FraudService.

func NewFraudService(detector fraud.FraudDetector, transactionServiceURL string) *FraudService {

	// Fallback to default URL if not provided

	if transactionServiceURL == "" {

		transactionServiceURL = os.Getenv("TRANSACTION_SERVICE_URL")

		if transactionServiceURL == "" {

			transactionServiceURL = "http://transaction-service:7000" // Default internal Docker service URL

		}

	}

	return &FraudService{detector: detector, transactionServiceURL: transactionServiceURL}

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

// This method will now be adapted for actual link validation.

func (s *FraudService) TrackPaymentLink(ctx context.Context, linkData map[string]interface{}) (bool, string, error) {

	urlStr, ok := linkData["url"].(string)

	if !ok || urlStr == "" {

		return true, "Invalid or missing URL in link data", nil // Treat as suspicious due to bad input

	}

	return s.ValidatePaymentLink(ctx, urlStr)

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

		Score: 0.1,

		Decision: "approve",
	}

}

func (s *FraudService) Rules(_ context.Context) []string {

	return []string{"velocity_check", "ip_blacklist"}

}

// GetTransactionDetailsByReference fetches transaction details from the transaction service by reference.

func (s *FraudService) GetTransactionDetailsByReference(ctx context.Context, reference string) (*TransactionResponse, error) {

	if s.transactionServiceURL == "" {

		return nil, fmt.Errorf("transaction service URL not configured")

	}

	url := fmt.Sprintf("%s/transactions/%s", strings.TrimRight(s.transactionServiceURL, "/"), reference)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)

	if err != nil {

		return nil, fmt.Errorf("failed to create request to transaction service: %w", err)

	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {

		return nil, fmt.Errorf("failed to get transaction from transaction service: %w", err)

	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {

		return nil, fmt.Errorf("transaction with reference %s not found", reference)

	}

	if resp.StatusCode != http.StatusOK {

		bodyBytes, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("transaction service returned status %d: %s", resp.StatusCode, string(bodyBytes))

	}

	var transactionResponse TransactionResponse

	if err := json.NewDecoder(resp.Body).Decode(&transactionResponse); err != nil {

		return nil, fmt.Errorf("failed to decode transaction response from transaction service: %w", err)

	}

	return &transactionResponse, nil

}

// ValidatePaymentLink parses a payment link URL and validates its parameters against transaction service.

func (s *FraudService) ValidatePaymentLink(ctx context.Context, linkUrl string) (bool, string, error) {

	parsedURL, err := url.Parse(linkUrl)

	if err != nil {

		return true, "Invalid payment link URL format", nil // Suspicious due to malformed URL

	}

	queryParams := parsedURL.Query()

	// Extract parameters from URL

	ref := queryParams.Get("ref")

	merchantIDStr := queryParams.Get("merchant_id")

	amountStr := queryParams.Get("amount")

			currency := queryParams.Get("currency")

			mode := queryParams.Get("mode") // For future use, if 'open' links have different validation

			_ = mode

	if ref == "" || merchantIDStr == "" || amountStr == "" || currency == "" {

		return true, "Missing required parameters in payment link (ref, merchant_id, amount, currency)", nil

	}

	// Convert merchantID and amount

	merchantID, err := strconv.Atoi(merchantIDStr)

	if err != nil {

		return true, "Invalid merchant_id format in payment link", nil

	}

	amount, err := strconv.ParseInt(amountStr, 10, 64)

	if err != nil {

		return true, "Invalid amount format in payment link", nil

	}

	// Fetch original transaction details

	originalTx, err := s.GetTransactionDetailsByReference(ctx, ref)

	if err != nil {

		// If transaction not found, it could be tampered or a non-existent link

		if strings.Contains(err.Error(), "not found") {

			return true, fmt.Sprintf("Transaction reference %s not found for payment link", ref), nil

		}

		return true, fmt.Sprintf("Error fetching original transaction details: %v", err), nil

	}

	// Validate against original transaction details

	if originalTx.MerchantID != merchantID {

		return true, fmt.Sprintf("Merchant ID mismatch: link has %d, original has %d", merchantID, originalTx.MerchantID), nil

	}

	if originalTx.Amount != amount {

		return true, fmt.Sprintf("Amount mismatch: link has %d, original has %d", amount, originalTx.Amount), nil

	}

	if originalTx.Currency != currency {

		return true, fmt.Sprintf("Currency mismatch: link has %s, original has %s", currency, originalTx.Currency), nil

	}

	// Additional checks for 'mode' or other parameters can be added here.

	// For example, if mode=open implies amount should be 0 in the original tx or be ignored.

	return false, "Payment link is legitimate", nil

}
