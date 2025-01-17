package auction

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

const (
	auctionInterval = 500 * time.Millisecond // Interval between auction processing steps.
)

type AuctionWorker struct {
	chainID int64
	state   *AuctionState // State of the auction managed by this worker.
	mu      sync.RWMutex  // Mutex for thread-safe access to the auction state.
}

// NewAuctionWorker creates a new auction worker for a specific chain.
func NewAuctionWorker(chainID int64) *AuctionWorker {
	return &AuctionWorker{
		chainID: chainID,
		state:   &AuctionState{},
	}
}

// StartAuction initializes a new auction with the given auction ID and information.
func (w *AuctionWorker) StartAuction(auction_info AuctionInfo) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.state.AuctionInfo = auction_info
	w.state.BidList = []Bid{}
	w.state.IsEnded = false

	fmt.Printf("[Worker %d] Starting new auction with ID: %s\n", w.chainID, auction_info.AuctionID)
	return nil
}

// SubmitBidList adds a batch of bids to the auction.
func (w *AuctionWorker) AddBids(auctionID string, BidList []Bid) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state.IsEnded {
		return errors.New("auction has already ended")
	}

	if w.state.AuctionInfo.AuctionID != auctionID {
		return errors.New("invalid auction ID")
	}

	w.state.BidList = append(w.state.BidList, BidList...)
	fmt.Printf("[Worker %d] Bid batch received: %d bids\n", w.chainID, len(BidList))
	return nil
}

// GetAuctionInfo returns the information of the current auction.
func (w *AuctionWorker) GetAuctionInfo() (AuctionInfo, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return w.state.AuctionInfo, nil
}

// GetLatestTob is a placeholder for retrieving the latest state of the order book.
func (w *AuctionWorker) GetLatestTob() ([]Tx, error) {
	return nil, nil
}

// GetAuctionState retrieves the current state of the auction.
func (w *AuctionWorker) GetAuctionState() (AuctionState, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()

	return *w.state, nil
}

// ProcessAuction runs a loop that processes the auction periodically.
// The loop terminates when the provided context is canceled.
func (w *AuctionWorker) ProcessAuction(ctx context.Context) {
	ticker := time.NewTicker(auctionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			err := w.processAuction()
			if err != nil {
				fmt.Printf("[Worker %d] Error processing auction: %v\n", w.chainID, err)
			}

			// Check if the auction has ended and terminate the loop.
			w.mu.RLock()
			if w.state.IsEnded {
				w.mu.RUnlock()
				fmt.Printf("[Worker %d] Auction has ended. Stopping ProcessAuction loop.\n", w.chainID)
				return
			}
			w.mu.RUnlock()

		case <-ctx.Done():
			fmt.Printf("[Worker %d] Stopping auction processing.\n", w.chainID)
			return
		}
	}
}

// processAuction performs the actual logic to manage the auction's state.
// It handles tasks such as starting, processing, and ending the auction.
func (w *AuctionWorker) processAuction() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state.IsEnded {
		return nil
	}

	now := time.Now().UnixMilli()
	startTime := w.state.AuctionInfo.StartTime.UnixMilli()
	endTime := w.state.AuctionInfo.EndTime.UnixMilli()

	// Auction has not started yet.
	if now < startTime {
		return nil
	}

	// Auction has ended.
	if now >= endTime {
		w.state.IsEnded = true

		// Sort BidList by BidAmount in descending order.
		sort.Slice(w.state.BidList, func(i, j int) bool {
			return w.state.BidList[i].BidAmount > w.state.BidList[j].BidAmount
		})

		// Extract all transactions from the sorted bids into SortedTxList.
		w.state.SortedTxList = []Tx{}
		for _, bid := range w.state.BidList {
			w.state.SortedTxList = append(w.state.SortedTxList, bid.TxList...)
		}

		// TODO: Get TD Quote and send to auction manager via gRPC channel.

		fmt.Printf("[Worker %d] Auction ended. Sorted transactions: %d\n", w.chainID, len(w.state.SortedTxList))

		// Clear the auction state after completion
		w.state = &AuctionState{}

		return nil
	}

	// If auction is still running, keep SortedTxList up to date.
	sort.Slice(w.state.BidList, func(i, j int) bool {
		return w.state.BidList[i].BidAmount > w.state.BidList[j].BidAmount
	})

	w.state.SortedTxList = []Tx{}
	for _, bid := range w.state.BidList {
		w.state.SortedTxList = append(w.state.SortedTxList, bid.TxList...)
	}

	fmt.Printf("[Worker %d] Auction processing. Current sorted transactions: %d\n", w.chainID, len(w.state.SortedTxList))
	return nil
}
