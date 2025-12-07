package repository

import (
	"context"
	"time"
)

// TransactionRecord represents a simplified transaction record for fraud checks.
type TransactionRecord struct {
	ID        string
	Amount    float64
	Currency  string
	Timestamp time.Time
	// Add more relevant fields like merchant ID, customer ID, IP, device ID, etc.
}

// FraudDataRepository defines the interface for accessing data required for fraud checks.
type FraudDataRepository interface {
	GetTransactionHistory(ctx context.Context, customerID string, lookback time.Duration) ([]TransactionRecord, error)
	GetIPData(ctx context.Context, ipAddress string) (map[string]interface{}, error) // Placeholder for IP-related data
	GetDeviceData(ctx context.Context, deviceID string) (map[string]interface{}, error) // Placeholder for device-related data
}

// InMemoryFraudDataRepository is a simple in-memory implementation of FraudDataRepository.
type InMemoryFraudDataRepository struct {
	transactions map[string][]TransactionRecord // customerID -> transactions
	ipData       map[string]map[string]interface{}
	deviceData   map[string]map[string]interface{}
}

// NewInMemoryFraudDataRepository creates a new in-memory repository.
func NewInMemoryFraudDataRepository() *InMemoryFraudDataRepository {
	return &InMemoryFraudDataRepository{
		transactions: make(map[string][]TransactionRecord),
		ipData:       make(map[string]map[string]interface{}),
		deviceData:   make(map[string]map[string]interface{}),
	}
}

// GetTransactionHistory retrieves transaction history for a given customer.
func (r *InMemoryFraudDataRepository) GetTransactionHistory(ctx context.Context, customerID string, lookback time.Duration) ([]TransactionRecord, error) {
	// Simulate fetching historical transactions.
	// In a real scenario, this would query a database.
	if txns, ok := r.transactions[customerID]; ok {
		var recentTxns []TransactionRecord
		cutoff := time.Now().Add(-lookback)
		for _, txn := range txns {
			if txn.Timestamp.After(cutoff) {
				recentTxns = append(recentTxns, txn)
			}
		}
		return recentTxns, nil
	}
	return []TransactionRecord{}, nil
}

// GetIPData retrieves data associated with an IP address.
func (r *InMemoryFraudDataRepository) GetIPData(ctx context.Context, ipAddress string) (map[string]interface{}, error) {
	// Simulate fetching IP data.
	if data, ok := r.ipData[ipAddress]; ok {
		return data, nil
	}
	return nil, nil
}

// GetDeviceData retrieves data associated with a device ID.
func (r *InMemoryFraudDataRepository) GetDeviceData(ctx context.Context, deviceID string) (map[string]interface{}, error) {
	// Simulate fetching device data.
	if data, ok := r.deviceData[deviceID]; ok {
		return data, nil
	}
	return nil, nil
}

// AddTransaction is a helper to populate the in-memory repository for testing.
func (r *InMemoryFraudDataRepository) AddTransaction(customerID string, txn TransactionRecord) {
	r.transactions[customerID] = append(r.transactions[customerID], txn)
}

// AddIPData is a helper to populate the in-memory repository for testing.
func (r *InMemoryFraudDataRepository) AddIPData(ipAddress string, data map[string]interface{}) {
	r.ipData[ipAddress] = data
}

// AddDeviceData is a helper to populate the in-memory repository for testing.
func (r *InMemoryFraudDataRepository) AddDeviceData(deviceID string, data map[string]interface{}) {
	r.deviceData[deviceID] = data
}
