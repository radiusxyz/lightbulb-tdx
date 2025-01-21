package auction

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

// Client represents a wrapper for AuctionServiceClient.
type Client struct {
	client auctionpb.AuctionServiceClient
	conn   *grpc.ClientConn
}

// NewClient initializes a new Client instance using grpc.NewClientConn.
func NewClient(serverAddr string) (*Client, error) {
	// Set up connection options.
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	// Use grpc.NewClientConn to establish the connection.
	conn, err := grpc.NewClient(serverAddr, opts...)
	if err != nil {
		return nil, err
	}

	return &Client{
		client: auctionpb.NewAuctionServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close gracefully closes the gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// AddAuction sends a request to start a new auction.
func (ac *Client) AddAuction(chainID int64, auctionID string, startTime, endTime time.Time, parameters string) {
	req := &auctionpb.AddAuctionRequest{
		AuctionInfo: &auctionpb.AuctionInfo{
			ChainId:    chainID,
			AuctionId:  auctionID,
			StartTime:  startTime.UnixMilli(),
			EndTime:    endTime.UnixMilli(),
		},
	}

	resp, err := ac.client.AddAuction(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to add auction: %v", err)
	}

	log.Printf("AddAuction Response: Success=%v, Message=%s", resp.Success, resp.Message)
}

// SubmitBids sends a batch of bids to the server.
func (ac *Client) SubmitBids(chainID int64, auctionID string, bids []*auctionpb.Bid) {
	req := &auctionpb.SubmitBidsRequest{
		ChainId:   chainID,
		AuctionId: auctionID,
		BidList:   bids,
	}

	resp, err := ac.client.SubmitBids(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to submit bids: %v", err)
	}

	log.Printf("SubmitBids Response: Success=%v, Message=%s", resp.Success, resp.Message)
}

// GetAuctionInfo retrieves detailed information about an auction.
func (ac *Client) GetAuctionInfo(chainID int64) {
	req := &auctionpb.GetAuctionInfoRequest{
		ChainId: chainID,
	}

	resp, err := ac.client.GetAuctionInfo(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get auction info: %v", err)
	}

	log.Printf("GetAuctionInfo Response: AuctionInfo=%v", resp.AuctionInfo)
}

// GetLatestTob retrieves the transaction list of the latest block.
func (ac *Client) GetLatestTob(chainID int64) {
	req := &auctionpb.GetLatestTobRequest{
		ChainId: chainID,
	}

	resp, err := ac.client.GetLatestTob(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get latest TOB: %v", err)
	}

	log.Printf("GetLatestTob Response: TxList=%v", resp.TxList)
}

// GetAuctionState retrieves the current state of an auction.
func (ac *Client) GetAuctionState(chainID int64) {
	req := &auctionpb.GetAuctionStateRequest{
		ChainId: chainID,
	}

	resp, err := ac.client.GetAuctionState(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to get auction state: %v", err)
	}

	log.Printf("GetAuctionState Response: State=%v", resp.State)
}
