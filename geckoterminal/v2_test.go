package geckoterminal

import (
	"errors"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestClient_SimpleTokenPrice(t *testing.T) {
	httpClient := &http.Client{
		Timeout: time.Second * 10,
	}

	c := NewClient(httpClient)

	// Test case 1: Valid network and token addresses
	network := "arbitrum"
	tokenAddresses := []string{"0x4a24B101728e07A52053c13FB4dB2BcF490CAbc3", "0x82af49447d8a07e3bd95bd0d56f35241523fbab1"}
	expectedPrices := map[string]float64{
		"0x123456789": 1.23,
		"0xabcdef123": 4.56,
	}

	prices, err := c.SimpleTokenPrice(network, tokenAddresses)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !reflect.DeepEqual(prices, expectedPrices) {
		t.Errorf("Unexpected prices. Got: %v, want: %v", prices, expectedPrices)
	}

	// Test case 2: Invalid network
	network = "invalid_network"
	tokenAddresses = []string{"0x123456789"}
	expectedError := errors.New("Invalid network")

	_, err = c.SimpleTokenPrice(network, tokenAddresses)
	if err == nil {
		t.Error("Expected an error, but got nil")
	} else if err.Error() != expectedError.Error() {
		t.Errorf("Unexpected error. Got: %v, want: %v", err, expectedError)
	}

	// Add more test cases as needed
}
