package attest

import (
	"fmt"

	"github.com/google/go-tdx-guest/client"
	tdxpb "github.com/google/go-tdx-guest/proto/tdx"
	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

// TDXClientWrapper is a wrapper around the tdxClient package.
type TDXClientWrapper struct{}

// NewTDXClientWrapper creates a new TDXClientWrapper.
func NewTDXClientWrapper() *TDXClientWrapper {
    return &TDXClientWrapper{}
}

// GetQuoteProvider wraps tdxClient.GetQuoteProvider().
func (w *TDXClientWrapper) GetQuoteProvider() (interface{}, error) {
	return client.GetQuoteProvider()
}

// GetQuote wraps tdxClient.GetQuote().
func (w *TDXClientWrapper) GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) {
	return client.GetQuote(provider, reportData)
}

// TDXClientInterface defines the methods required for a TDX client.
type TDXClientInterface interface {
    GetQuoteProvider() (interface{}, error) // Returns a quote provider (interface to allow flexibility)
    GetQuote(provider interface{}, reportData [64]byte) (interface{}, error) // Returns a quote object
}

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

    return ConvertQuoteV4ToQuote(convertedQuote), nil
}