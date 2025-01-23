package test

import (
	"fmt"
	"log"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/radiusxyz/lightbulb-tdx/auction"
	"gopkg.in/yaml.v3"

	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

// Scenario defines the load test parameters.
type Scenario struct {
    AuctionTime      float64 `yaml:"auction_time"`
    OrderingInterval float64 `yaml:"ordering_interval"`
    ChainNum         int     `yaml:"chain_num"`
    AuctionNum       int     `yaml:"auction_num"`
    ClientNum        int     `yaml:"client_num"`
    RequestFreq      int     `yaml:"request_freq"`
}

// AuctionMeta holds each auction's info.
type AuctionMeta struct {
    ChainID   int64
    AuctionID string
    StartTime time.Time
    EndTime   time.Time
}

// BenchmarkAuctionWorker shows how to assign subsets of auctions to different clients.
func BenchmarkAuctionWorker(t *testing.B) {
    // 1) Load environment
    err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
    port := os.Getenv("PORT")
    serverAddr := os.Getenv("SERVER_ADDRESS")

    // 2) Load scenario
    scenarioData, err := os.ReadFile("scenario.yaml")
    if err != nil {
        t.Fatalf("Failed to read YAML file: %v", err)
    }
    var scenario Scenario
    if err := yaml.Unmarshal(scenarioData, &scenario); err != nil {
        t.Fatalf("Failed to unmarshal YAML: %v", err)
    }

    // 3) Prepare WaitGroup
    var wg sync.WaitGroup

    // 4) Manager Client (creates auctions)
    managerClient, err := auction.NewClient(fmt.Sprintf("%s:%s", serverAddr, port))
    if err != nil {
        t.Fatalf("Failed to create Manager Client: %v", err)
    }
    defer managerClient.Close()

    // Auctions: map[chainIndex] -> []AuctionMeta
    Auctions := make(map[int][]AuctionMeta, scenario.ChainNum)
    var auctionsMu sync.Mutex

    wg.Add(1)
    go func() {
        defer wg.Done()
        log.Println("[Manager Client] Adding auctions...")

        baseTime := time.Now().Add(time.Second)
        auctionDuration := time.Duration(scenario.AuctionTime * float64(time.Second))

        for c := 0; c < scenario.ChainNum; c++ {
            chainID := int64(c)
            for a := 0; a < scenario.AuctionNum; a++ {
                auctionID := fmt.Sprintf("chain_%d_auction_%d", c, a)

                startTime := baseTime.Add(time.Duration(a) * auctionDuration)
                endTime := startTime.Add(auctionDuration)

                // gRPC call to add auction
                managerClient.AddAuction(chainID, auctionID, startTime, endTime, "")

                // Store in local map
                auctionsMu.Lock()
                Auctions[c] = append(Auctions[c], AuctionMeta{
                    ChainID:   chainID,
                    AuctionID: auctionID,
                    StartTime: startTime,
                    EndTime:   endTime,
                })
                auctionsMu.Unlock()

                log.Printf("[Manager] Auction '%s' on chain %d added (start=%v, end=%v)",
                    auctionID, chainID, startTime.Format("15:04:05.000"), endTime.Format("15:04:05.000"))
            }
        }
        log.Println("[Manager Client] Finished adding auctions.")
    }()

    // 5) General Clients: each client only sees the auctions for chain (clientID % chainNum)
    //    We'll wait for manager to finish. Or we can do a simple Sleep to allow manager to finish first.
    //    In a real scenario, consider a Condition Variable or other sync method.
    wg.Wait()

    // Now manager is done => we can safely read from Auctions
    // Start clients
    for i := 0; i < scenario.ClientNum; i++ {
        wg.Add(1)
        go func(clientID int) {
            defer wg.Done()

            // Make a new gRPC client
            client, err := auction.NewClient(fmt.Sprintf("%s:%s", serverAddr, port))
            if err != nil {
                log.Fatalf("[Client %d] Failed to create client: %v", clientID, err)
            }
            defer client.Close()

            // The chain index this client will handle
            chainIdx := clientID % scenario.ChainNum

            // Copy local auctions (for read-only usage)
            localAuctions := make([]AuctionMeta, len(Auctions[chainIdx]))
            copy(localAuctions, Auctions[chainIdx])

            // Set up request parameters
            requestInterval := time.Duration(1e9 / scenario.RequestFreq)

            log.Printf("[Client %d] Handling chain %d with %d auctions", clientID, chainIdx, len(localAuctions))

            // Send bids only for that chain's auctions
            for _, auction := range localAuctions {
                for bid := 0; ; bid++ {
                    if time.Now().Before(auction.StartTime) {
                        time.Sleep(requestInterval)
                        continue
                    }
                    if time.Now().After(auction.EndTime) {
                        break
                    }
                    bids := []*auctionpb.Bid{
                        {
                            BidAmount: int64(bid + 1),
                            TxList: []*auctionpb.Tx{
                                {TxData: fmt.Sprintf("tx-client%d-%d", clientID, bid)},
                                {TxData: fmt.Sprintf("tx-client%d-%d", clientID, bid+1)},
                            },
                        },
                    }
                    client.SubmitBids(auction.ChainID, auction.AuctionID, bids)
                    time.Sleep(requestInterval)
                }
            }

            log.Printf("[Client %d] Finished all requests.", clientID)
        }(i)
    }

    // 6) Wait for clients to finish
    wg.Wait()
    log.Println("[Benchmark] All clients done.")
}