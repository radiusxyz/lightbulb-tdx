package tdx

import (
	"fmt"
	"os"

	"github.com/google/go-tdx-guest/client"

	"github.com/radiusxyz/lightbulb-tdx/utils"

	tdxpb "github.com/google/go-tdx-guest/proto/tdx"
	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

// TDXClientInterface defines the methods required for a TDX client.
type TDXClientInterface interface {
	GetQuoteProvider() (interface{}, error)                                  // Returns a quote provider
	GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) // Returns a quote object
	GetRtmr() ([]byte, error)                                                // Returns the RTMR value
}

// TDXClient is a wrapper around the tdxClient package.
type TDXClient struct{
	rtmrProvider *RtmrProvider
}

// NewTDXClient creates a new TDXClient.
func NewTDXClient() *TDXClient {
	return &TDXClient{
		rtmrProvider: DefaultRtmrProvider(),
	}
}

// GetQuoteProvider wraps tdxClient.GetQuoteProvider().
func (w *TDXClient) GetQuoteProvider() (interface{}, error) {
	return client.GetQuoteProvider()
}

// GetQuote wraps tdxClient.GetQuote().
func (w *TDXClient) GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) {
	return client.GetQuote(provider, reportData)
}

type MockTDXClient struct {
	rtmrProvider *RtmrProvider
}

func NewMockTDXClient() *MockTDXClient {
	return &MockTDXClient{
		rtmrProvider: DefaultRtmrProvider(),
	}
}

type MockQuoteProvider struct {}

func (m *MockTDXClient) GetQuoteProvider() (interface{}, error) {
	return &MockQuoteProvider{}, nil
}

func (m *MockTDXClient) GetQuote(_provider interface{}, reportData [64]byte) (interface{}, error) {
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
func GetQuote(tdxClient TDXClientInterface) (*attestpb.Quote, error) {
	// Get the quote provider
	quoteProvider, err := tdxClient.GetQuoteProvider()
	if err != nil {
		return nil, fmt.Errorf("failed to get quote provider: %w", err)
	}

	reportDataInput := make([]byte, 64)

	if os.Getenv("TDX_VERSION") == "1.0" {
		reportDataInput, err = tdxClient.GetRtmr()
		if err != nil {
			return nil, fmt.Errorf("failed to get RTMR: %w", err)
		}
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

func (c *TDXClient) GetRtmr() ([]byte, error) {
	// Update RTMR[2] with IMA Event Logs
	err := c.rtmrProvider.UpdateImaRtmr()
	if err != nil {
		return nil, fmt.Errorf("failed to update IMA RTMR: %w", err)
	}

	// Get the RTMR[2] value
	return c.rtmrProvider.GetRtmrValues()[2], nil;
}

func (c *MockTDXClient) GetRtmr() ([]byte, error) {
	// Update RTMR[2] with IMA Event Logs
	err := c.rtmrProvider.UpdateImaRtmr()
	if err != nil {
		return nil, fmt.Errorf("failed to update IMA RTMR: %w", err)
	}

	// Get the RTMR[2] value
	return c.rtmrProvider.GetRtmrValues()[2], nil;
}