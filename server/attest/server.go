package attest

import (
	"context"
	"log"

	tdxClient "github.com/google/go-tdx-guest/client"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
	tdxpb "github.com/radiusxyz/lightbulb-tdx/proto/tdx"
)

type Server struct {
	attestpb.UnimplementedAttestServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// Attest implements the AttestServiceServer interface.
func (s *Server) GetQuote(ctx context.Context, req *attestpb.GetQuoteRequest) (*attestpb.GetQuoteResponse, error) {
	log.Printf("Received Attest request for report_data=%x", req.GetReportData())

	quoteProvider, err := tdxClient.GetQuoteProvider()
	if err != nil {
		log.Fatalf("Failed to get quote provider: %v", err)
	}
	var reportData [64]byte
	copy(reportData[:], req.GetReportData())

	quote, err := tdxClient.GetQuote(quoteProvider, reportData)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	quoteProto := quote.(*tdxpb.QuoteV4)

	// Debug: Print RTMR values.
	for i, rtmr := range quoteProto.TdQuoteBody.Rtmrs {
		log.Printf("rtmr[%d]: %x", i, rtmr)
	}

	return &attestpb.GetQuoteResponse{
		Quote: quoteProto,
	}, nil
}
