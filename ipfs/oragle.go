package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// SignRequest represents the request body for the oracle API
type SignRequest struct {
	CID string `json:"cid"`
}

// SignatureResponse represents a single signature in the response
type SignatureResponse struct {
	Signer    string `json:"signer"`
	Signature string `json:"signature"`
}

// OracleClient defines the interface for interacting with the IPFS oracle
type OracleClient interface {
	GetSignaturesForCID(ctx context.Context, cid string) ([]SignatureResponse, error)
}

// HTTPOracleClient implements the OracleClient interface using HTTP
type HTTPOracleClient struct {
	OracleURL      string
	DefaultTimeout time.Duration
	client         *http.Client
}

// NewHTTPOracleClient creates a new HTTP oracle client
func NewHTTPOracleClient(oracleURL string, timeout time.Duration) *HTTPOracleClient {
	return &HTTPOracleClient{
		OracleURL:      oracleURL,
		DefaultTimeout: timeout,
		client:         &http.Client{Timeout: timeout},
	}
}

// GetSignaturesForCID calls the IPFS oracle API to get signatures for a pinned CID
func (c *HTTPOracleClient) GetSignaturesForCID(ctx context.Context, cid string) ([]SignatureResponse, error) {
	// Create request body
	reqBody, err := json.Marshal(SignRequest{CID: cid})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/sign", c.OracleURL), bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send request using the shared client
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("oracle API returned non-OK status: %d, body: %s", resp.StatusCode, string(body))
	}

	// Parse response
	var signatures []SignatureResponse
	if err := json.NewDecoder(resp.Body).Decode(&signatures); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return signatures, nil
}

// MockOracleClient implements the OracleClient interface for testing
type MockOracleClient struct {
	// Predefined responses for specific CIDs
	MockResponses map[string][]SignatureResponse
	// Default response when CID is not found in MockResponses
	DefaultResponse []SignatureResponse
}

// NewMockOracleClient creates a new mock oracle client
func NewMockOracleClient() *MockOracleClient {
	return &MockOracleClient{
		MockResponses: make(map[string][]SignatureResponse),
		DefaultResponse: []SignatureResponse{
			{
				Signer:    "0x0",
				Signature: "0xB00BA",
			},
		},
	}
}

// AddMockResponse adds a mock response for a specific CID
func (m *MockOracleClient) AddMockResponse(cid string, responses []SignatureResponse) {
	m.MockResponses[cid] = responses
}

// GetSignaturesForCID returns mock signatures for a CID
func (m *MockOracleClient) GetSignaturesForCID(ctx context.Context, cid string) ([]SignatureResponse, error) {
	// Check for context cancellation
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Return predefined response if available, otherwise return default
	if responses, ok := m.MockResponses[cid]; ok {
		return responses, nil
	}
	return m.DefaultResponse, nil
}
