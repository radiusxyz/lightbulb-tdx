package auction

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/radiusxyz/lightbulb-tdx/tdx"
)

const (
	auctionInterval = 500 * time.Millisecond // Interval for processing auction state.
)

// AuctionWorker manages auctions in a queue, ensuring they are processed by start time.
type AuctionWorker struct {
	chainID      int64                       // Unique identifier for the worker.
	state        *AuctionState               // Holds the current state of the auction being processed.
	mu           sync.RWMutex                // RWMutex for protecting shared resources.
	queueCond    *sync.Cond                  // Condition variable to handle empty queue waiting.
	auctionQueue []AuctionInfo               // Queue of auctions sorted by StartTime.
	interruptCh  chan struct{}   		     // Channel to interrupt waiting when queue changes.
	tdxClient    tdx.TDXClientInterface      // TDX client for quote generation.
}

// NewAuctionWorker initializes a new AuctionWorker and starts its queue processor.
func NewAuctionWorker(chainID int64) *AuctionWorker {
	var tdxClient tdx.TDXClientInterface
	
	env := os.Getenv("ENV")

	if env == "TDX" {
		tdxClient = tdx.NewTDXClient()
	} else if env == "MOCK_TDX" {
		tdxClient = tdx.NewMockTDXClient()
	} else {
		log.Printf("[Warning] Unknown environment '%s'. Defaulting to MockTDXClient.", env)
		tdxClient = tdx.NewMockTDXClient()
	}

	worker := &AuctionWorker{
		chainID:      chainID,
		state:        &AuctionState{},
		tdxClient:    tdxClient,
		interruptCh:  make(chan struct{}, 1),
	}
	worker.queueCond = sync.NewCond(&worker.mu)

	// Start queue processing in a separate goroutine.
	go func() {
		ctx := context.Background()
		worker.StartQueueProcessor(ctx)
	}()

	return worker
}

// initializeAuction sets up the auction state before it starts.
func (w *AuctionWorker) initializeAuction(info AuctionInfo) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.state.AuctionInfo = info
	w.state.BidList = []Bid{}
	w.state.IsEnded = false

	log.Printf("[Worker %d] Initializing auction (ID: %s)\n", w.chainID, info.AuctionID)
	return nil
}

// ProcessAuction periodically updates the state of the current auction until it ends or is canceled.
func (w *AuctionWorker) ProcessAuction(ctx context.Context) {
	ticker := time.NewTicker(auctionInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if w.isEnded() {
				log.Printf("[Worker %d] Auction has ended. Stopping ProcessAuction.\n", w.chainID)
				return
			}
			err := w.processAuctionLogic()
			if err != nil {
				log.Printf("[Worker %d] Auction error: %v\n", w.chainID, err)
			}
		case <-ctx.Done():
			log.Printf("[Worker %d] Context canceled. Stopping ProcessAuction.\n", w.chainID)
			return
		}
	}
}

// processAuctionLogic updates the auction state and processes bids.
func (w *AuctionWorker) processAuctionLogic() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state.IsEnded {
		return nil
	}
	now := time.Now().UnixMilli()
	start := w.state.AuctionInfo.StartTime.UnixMilli()
	end := w.state.AuctionInfo.EndTime.UnixMilli()

	if now < start {
		return nil // Auction has not started yet.
	}
	if now >= end {
		// Auction has ended.
		w.state.IsEnded = true
		sort.Slice(w.state.BidList, func(i, j int) bool {
			return w.state.BidList[i].BidAmount > w.state.BidList[j].BidAmount
		})
		w.state.SortedTxList = nil
		for _, bid := range w.state.BidList {
			w.state.SortedTxList = append(w.state.SortedTxList, bid.TxList...)
		}
		log.Printf("[Worker %d] Auction ended with %d transactions.\n", w.chainID, len(w.state.SortedTxList))
		return nil
	}

	// Process ongoing bids.
	sort.Slice(w.state.BidList, func(i, j int) bool {
		return w.state.BidList[i].BidAmount > w.state.BidList[j].BidAmount
	})
	w.state.SortedTxList = nil
	for _, bid := range w.state.BidList {
		w.state.SortedTxList = append(w.state.SortedTxList, bid.TxList...)
	}

	log.Printf("[Worker %d] Auction running. Sorted transactions: %d\n", w.chainID, len(w.state.SortedTxList))
	return nil
}

// isEnded checks if the current auction has ended.
func (w *AuctionWorker) isEnded() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state.IsEnded
}

// AddBids adds new bids to the current auction.
func (w *AuctionWorker) AddBids(auctionID string, bids []Bid) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.state.IsEnded {
		return errors.New("auction has already ended")
	}
	if w.state.AuctionInfo.AuctionID != auctionID {
		return errors.New("invalid auction ID")
	}
	w.state.BidList = append(w.state.BidList, bids...)
	log.Printf("[Worker %d] Received %d bids\n", w.chainID, len(bids))
	return nil
}

// AddAuction adds a new auction to the queue and interrupts waiting if necessary.
func (w *AuctionWorker) AddAuction(info AuctionInfo) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	now := time.Now()
	if info.StartTime.Before(now) {
		return fmt.Errorf("invalid auction start time: %s is before now %s", info.StartTime, now)
	}
	if info.EndTime.Before(info.StartTime) {
		return fmt.Errorf("end time %s is before start time %s", info.EndTime, info.StartTime)
	}
	for _, a := range w.auctionQueue {
		if a.AuctionID == info.AuctionID {
			return fmt.Errorf("auction ID %s already exists", info.AuctionID)
		}
	}
	w.auctionQueue = append(w.auctionQueue, info)
	sort.Slice(w.auctionQueue, func(i, j int) bool {
		return w.auctionQueue[i].StartTime.Before(w.auctionQueue[j].StartTime)
	})
	log.Printf("[Worker %d] Enqueued auction (ID: %s)\n", w.chainID, info.AuctionID)

	w.queueCond.Signal()
	w.interrupt()

	return nil
}

// StartQueueProcessor processes the auction queue in order of StartTime.
func (w *AuctionWorker) StartQueueProcessor(ctx context.Context) {
	for {
		w.mu.Lock()
		for len(w.auctionQueue) == 0 {
			if ctx.Err() != nil {
				w.mu.Unlock()
				return
			}
			w.queueCond.Wait()
			if ctx.Err() != nil {
				w.mu.Unlock()
				return
			}
		}

		nextAuction := w.auctionQueue[0]
		now := time.Now()

		if now.Before(nextAuction.StartTime) {
			waitDuration := nextAuction.StartTime.Sub(now)
			w.mu.Unlock()
			select {
			case <-time.After(waitDuration):
			case <-w.interruptCh:
			case <-ctx.Done():
				return
			}
			continue
		}

		w.auctionQueue = w.auctionQueue[1:]
		w.mu.Unlock()
		w.runAuction(ctx, nextAuction)
	}
}

// runAuction initializes and processes a single auction.
func (w *AuctionWorker) runAuction(ctx context.Context, info AuctionInfo) {
	if err := w.initializeAuction(info); err != nil {
		log.Printf("[Worker %d] Failed to start auction %s: %v\n", w.chainID, info.AuctionID, err)
		return
	}

	// Defer getting a quote after the auction ends.
	defer func() {
		quote, err := tdx.GetQuote(w.tdxClient)
		if err != nil {
			log.Printf("[Worker %d] Failed to get quote: %v\n", w.chainID, err)
		} else {
			log.Printf("[Worker %d] Quote: %v\n", w.chainID, quote)
		}
	}()

	subCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	done := make(chan struct{})
	go func() {
		w.ProcessAuction(subCtx)
		close(done)
	}()

	select {
	case <-done:
		log.Printf("[Worker %d] Auction (ID: %s) completed.\n", w.chainID, info.AuctionID)
	case <-ctx.Done():
		log.Printf("[Worker %d] Context canceled. Stopping auction (ID: %s).\n", w.chainID, info.AuctionID)
	}
}

// GetAuctionInfo retrieves the current auction info.
func (w *AuctionWorker) GetAuctionInfo() AuctionInfo {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.state.AuctionInfo
}

// GetAuctionState retrieves the current auction state.
func (w *AuctionWorker) GetAuctionState() AuctionState {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return *w.state
}

// GetLatestTob retrieves the latest order book transactions (placeholder).
func (w *AuctionWorker) GetLatestTob() ([]Tx, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return nil, nil
}

// interrupt triggers an immediate queue check by interrupting any ongoing wait.
func (w *AuctionWorker) interrupt() {
	select {
	case w.interruptCh <- struct{}{}:
	default:
	}
}