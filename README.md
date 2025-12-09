# Fraud Service

Fraud detection and prevention service for KodraPay platform.

## Overview

The Fraud Service provides real-time fraud detection capabilities through rule-based analysis and transaction monitoring. It helps identify suspicious transactions before they are processed.

## Features

- **Rule-based fraud detection** - Configurable fraud rules (high amounts, suspicious IPs, high velocity)
- **Transaction checking** - Real-time validation of transactions
- **Payment link validation** - Fraud checks for payment link transactions
- **Payment channel validation** - Verification of payment channels
- **API key authentication** - Secure access control via X-API-Key header
- **Rate limiting** - 5 requests per second per API key

## Configuration

### Environment Variables

Create a `.env` file based on `.env.example`:

```bash
PORT=7012
POSTGRES_URL=postgres://user:password@localhost:5432/kodrapay
REDIS_URL=redis://localhost:6379
API_KEY=your-secure-api-key-here
```

### API Key Setup

**IMPORTANT:** The fraud service requires an API key for authentication.

1. **Generate a secure API key:**
   ```bash
   # Use openssl to generate a random key
   openssl rand -hex 32
   ```

2. **Set the API key in your environment:**
   ```bash
   export API_KEY=your-generated-key
   ```

3. **For Docker deployment:** The API key is already configured in `docker-compose.yml`:
   ```yaml
   environment:
     API_KEY: my-secret-api-key
   ```

   **For production, replace with a secure key!**

4. **Include the API key in all requests:**
   ```bash
   curl -X POST http://localhost:7012/fraud/check-transaction \
     -H "X-API-Key: your-api-key" \
     -H "Content-Type: application/json" \
     -d '{"transaction_id": "tx_123", "amount": 100000, "currency": "NGN"}'
   ```

## API Endpoints

### Health Check
```
GET /health
```
No authentication required.

### Check Transaction
```
POST /fraud/check-transaction
Headers: X-API-Key: <your-api-key>
Body: {
  "transaction_id": "string",
  "merchant_id": "string",
  "amount": number,
  "currency": "string",
  "customer_ip": "string"
}
```

### Track Payment Link
```
POST /fraud/track-payment-link
Headers: X-API-Key: <your-api-key>
Body: {
  "payment_link_id": "string",
  "merchant_id": "string",
  "amount": number,
  "customer_ip": "string"
}
```

### Validate Payment Channel
```
POST /fraud/validate-payment-channel
Headers: X-API-Key: <your-api-key>
Body: {
  "channel_id": "string",
  "merchant_id": "string"
}
```

## Fraud Rules

The service supports the following fraud detection rules:

- **HIGH_AMOUNT_TRANSACTION** - Flags transactions above 10,000,000 (in smallest currency unit)
- **SUSPICIOUS_IP_ORIGIN** - Checks IP addresses against known malicious patterns
- **HIGH_VELOCITY_CUSTOMER** - Detects unusual transaction frequency per merchant

Rules can be configured with:
- `enabled`: true/false
- `severity`: low/medium/high/critical
- `action`: allow/review/block

## Running the Service

### Local Development
```bash
cd fraud-service
go run ./cmd/fraud-service
```

### Docker
```bash
docker-compose up fraud-service
```

### Build
```bash
go build -o bin/fraud-service ./cmd/fraud-service
./bin/fraud-service
```

## Rate Limiting

- Rate limit: 5 requests per second per API key
- Exceeding the limit returns HTTP 429 (Too Many Requests)

## Security Considerations

1. **Never commit your API key to version control**
2. **Use environment variables** for API key configuration
3. **Generate unique keys** for each environment (dev, staging, production)
4. **Rotate API keys** periodically
5. **Use HTTPS** in production to protect API keys in transit

## Development

### Project Structure
```
fraud-service/
├── cmd/fraud-service/main.go        # Entry point
├── internal/
│   ├── config/                      # Configuration loading
│   ├── middleware/                  # API key auth, rate limiting
│   ├── handlers/                    # HTTP request handlers
│   ├── services/                    # Business logic
│   ├── repositories/                # Data access (in-memory)
│   ├── models/                      # Domain models
│   └── routes/                      # Route registration
├── .env.example                     # Environment template
├── Dockerfile                       # Container definition
└── README.md                        # This file
```

## Troubleshooting

### "Unauthorized" error
- Check that you're including the `X-API-Key` header
- Verify the API key matches the one configured in your environment
- Ensure the API key doesn't have extra whitespace

### Rate limit exceeded
- The service limits requests to 5 per second per API key
- Implement exponential backoff in your client
- Consider requesting a higher rate limit for your use case

## Future Enhancements

- [ ] Machine learning-based fraud detection
- [ ] Integration with third-party fraud detection services
- [ ] Persistent storage for fraud data
- [ ] Advanced IP geolocation checking
- [ ] Device fingerprinting
- [ ] Behavioral analysis
- [ ] Real-time alerting and notifications
