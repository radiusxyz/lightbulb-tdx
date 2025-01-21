package auction

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/radiusxyz/lightbulb-tdx/attest"

	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

type Server struct {
	auctionpb.UnimplementedAuctionServiceServer

	workers      map[int64]*AuctionWorker    // Workers mapped by chain ID
	rtmrExtender *attest.IMAEventLogExtender // Extender for RTMR values
	mu           sync.RWMutex                // Mutex to ensure thread-safe access to the workers map.
}

// NewServer initializes a new gRPC server instance.
func NewServer() *Server {
	return &Server{
		workers: make(map[int64]*AuctionWorker),
		rtmrExtender: attest.DefaultIMAEventLogExtender(),
	}
}

// AddAuction handles the gRPC call to start an auction.
func (s *Server) AddAuction(ctx context.Context, req *auctionpb.AddAuctionRequest) (*auctionpb.AddAuctionResponse, error) {
	pbInfo := req.GetAuctionInfo()
	info := ConvertProtobufAuctionInfoToDomain(pbInfo)

	s.mu.Lock()
	defer s.mu.Unlock()

	// Retrieve or create the worker for the chain
	worker, exists := s.workers[info.ChainID]
	if !exists {
		worker = NewAuctionWorker(info.ChainID, s.rtmrExtender)
		s.workers[info.ChainID] = worker
	}

	// Debug: Print the worker map
	for chainID, worker := range s.workers {
		log.Printf("Worker %d: %v\n", chainID, worker)
	}

	// Start the auction
	err := worker.AddAuction(info)
	if err != nil {
		return &auctionpb.AddAuctionResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &auctionpb.AddAuctionResponse{
		Success: true,
		Message: "Auction started successfully",
	}, nil
}

// SubmitBids handles the gRPC call to submit a batch of bids.
func (s *Server) SubmitBids(ctx context.Context, req *auctionpb.SubmitBidsRequest) (*auctionpb.SubmitBidsResponse, error) {
	chainID := req.GetChainId()
	auctionID := req.GetAuctionId()
	pbBidList := req.GetBidList()

	s.mu.RLock()
	worker, exists := s.workers[chainID]
	s.mu.RUnlock()

	if !exists {
		return &auctionpb.SubmitBidsResponse{
			Success: false,
			Message: "Chain not found",
		}, nil
	}

	bidList := ConvertProtobufBidsToDomain(pbBidList)

	err := worker.AddBids(auctionID, bidList)
	if err != nil {
		return &auctionpb.SubmitBidsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &auctionpb.SubmitBidsResponse{
		Success: true,
		Message: "Bids submitted successfully",
	}, nil
}

// GetAuctionInfo retrieves detailed information about a specific auction.
func (s *Server) GetAuctionInfo(ctx context.Context, req *auctionpb.GetAuctionInfoRequest) (*auctionpb.GetAuctionInfoResponse, error) {
	chainID := req.GetChainId()

	s.mu.RLock()
	worker, exists := s.workers[chainID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("chain not found")
	}

	info := worker.GetAuctionInfo()

	return &auctionpb.GetAuctionInfoResponse{
		AuctionInfo: ConvertDomainAuctionInfoToProtobuf(info),
	}, nil
}

// GetLatestTob retrieves the Tx list of the latest block.
func (s *Server) GetLatestTob(ctx context.Context, req *auctionpb.GetLatestTobRequest) (*auctionpb.GetLatestTobResponse, error) {
	chainID := req.GetChainId()

	s.mu.RLock()
	worker, exists := s.workers[chainID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("chain not found")
	}

	txList, err := worker.GetLatestTob()
	if err != nil {
		return nil, err
	}

	return &auctionpb.GetLatestTobResponse{
		TxList: ConvertDomainTxsToProtobuf(txList),
	}, nil
}

// GetAuctionState retrieves the current state of an auction.
func (s *Server) GetAuctionState(ctx context.Context, req *auctionpb.GetAuctionStateRequest) (*auctionpb.GetAuctionStateResponse, error) {
	chainID := req.GetChainId()

	s.mu.RLock()
	worker, exists := s.workers[chainID]
	s.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("chain not found")
	}

	state := worker.GetAuctionState()

	return &auctionpb.GetAuctionStateResponse{
		State: ConvertDomainAuctionStateToProtobuf(state),
	}, nil
}
