package attest

import (
	"context"
	"fmt"
	"log"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
)

type Server struct {
    attestpb.UnimplementedAttestServiceServer
    tdxClient TDXClientInterface
}

// NewServer creates a new server with a TDXClientInterface.
func NewServer(client TDXClientInterface) *Server {
    return &Server{
        tdxClient: client,
    }
}

func (s *Server) GetQuote(ctx context.Context, req *attestpb.GetQuoteRequest) (*attestpb.GetQuoteResponse, error) {
	// Get the quote
	quoteProto, err := GetQuote(req.GetReportData(), s.tdxClient)
	if err != nil {
		return nil, fmt.Errorf("failed to get quote: %v", err)
	}

    // Debug: Print RTMR values
    for i, rtmr := range quoteProto.TdQuoteBody.Rtmrs {
        log.Printf("rtmr[%d]: %x", i, rtmr)
    }

    return &attestpb.GetQuoteResponse{
        Quote: quoteProto,
    }, nil
}