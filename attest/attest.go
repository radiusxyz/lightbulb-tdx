package attest

import (
	"context"
	"encoding/json"
	"log"

	tdxClient "github.com/google/go-tdx-guest/client"
	tdxProto "github.com/google/go-tdx-guest/proto/tdx"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

type Server struct {
	attestpb.UnimplementedAttestServiceServer
}

func NewServer() *Server {
	return &Server{}
}

// Attest implements the AttestServiceServer interface.
func (s *Server) Attest(ctx context.Context, req *attestpb.AttestRequest) (*attestpb.AttestReply, error) {
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

	quoteProto := quote.(*tdxProto.QuoteV4)

	// Debug: Print RTMR values.
	for i, rtmr := range quoteProto.TdQuoteBody.Rtmrs {
		log.Printf("rtmr[%d]: %x", i, rtmr)
	}

	quoteBytes, err := json.MarshalIndent(quote, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal quote to JSON: %v", err)
	}

	return &attestpb.AttestReply{
		Quote: string(quoteBytes),
	}, nil
}
