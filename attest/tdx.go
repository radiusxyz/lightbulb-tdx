package attest

import (
	"fmt"

	"github.com/google/go-tdx-guest/client"

	"github.com/radiusxyz/lightbulb-tdx/utils"

	tdxpb "github.com/google/go-tdx-guest/proto/tdx"
	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

// TDXClientInterface defines the methods required for a TDX client.
type TDXClientInterface interface {
	GetQuoteProvider() (interface{}, error)                                  // Returns a quote provider
	GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) // Returns a quote object
}

// TDXClient is a wrapper around the tdxClient package.
type TDXClient struct{}

// NewTDXClient creates a new TDXClient.
func NewTDXClient() *TDXClient {
	return &TDXClient{}
}

// GetQuoteProvider wraps tdxClient.GetQuoteProvider().
func (w *TDXClient) GetQuoteProvider() (interface{}, error) {
	return client.GetQuoteProvider()
}

// GetQuote wraps tdxClient.GetQuote().
func (w *TDXClient) GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) {
	return client.GetQuote(provider, reportData)
}

// MockTDXClient is a mock implementation of TDXClientInterface for non-TDX environments.
type MockTDXClient struct{}

// NewMockTDXClient creates a new instance of MockTDXClient.
func NewMockTDXClient() *MockTDXClient {
	return &MockTDXClient{}
}

// GetQuoteProvider mocks the retrieval of a quote provider.
func (m *MockTDXClient) GetQuoteProvider() (interface{}, error) {
	// Return a dummy provider
	return "mockQuoteProvider", nil
}

// GetQuote mocks the retrieval of a quote.
func (m *MockTDXClient) GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) {
	if provider != "mockQuoteProvider" {
		return nil, fmt.Errorf("invalid mock provider: %v", provider)
	}

	// Create a mock quote
	mockQuote := &tdxpb.QuoteV4{
        Header: &tdxpb.Header{
            Version: 1,
        },
		TdQuoteBody: &tdxpb.TDQuoteBody{
            ReportData: reportData[:],
	    },
    }

	return mockQuote, nil
}

// GetQuote retrieves a TDX quote using the given TDX client implementation.
func GetQuote(reportDataInput []byte, tdxClient TDXClientInterface) (*attestpb.Quote, error) {
	// Check input validity
	if len(reportDataInput) > 64 {
		return nil, fmt.Errorf("report_data exceeds maximum length of 64 bytes")
	}

	// Get the quote provider
	quoteProvider, err := tdxClient.GetQuoteProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to get quote provider: %w", err)
	}

	// Prepare reportData array
	var reportData [64]byte
	copy(reportData[:], reportDataInput)

	// Get the quote
	quote, err := tdxClient.GetQuote(quoteProvider, reportData)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %w", err)
	}

	// Convert and return the Quote
	convertedQuote, ok := quote.(*tdxpb.QuoteV4)
	if !ok {
		return nil, fmt.Errorf("unexpected quote type: %T", quote)
	}

	return utils.ConvertQuoteV4ToQuote(convertedQuote), nil
}