package benchmark

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	benchpb "github.com/radiusxyz/lightbulb-tdx/proto/benchmark"
)

// Client represents a wrapper for AuctionServiceClient.
type Client struct {
	client benchpb.BenchmarkServiceClient
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
		client: benchpb.NewBenchmarkServiceClient(conn),
		conn:   conn,
	}, nil
}

// Close gracefully closes the gRPC connection.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Hello sends a hello request to the server.
func (c *Client) Hello() (*benchpb.HelloResponse, error) {
	return c.client.Hello(context.Background(), &benchpb.HelloRequest{})
}

// CPUIntensive sends a CPU-intensive computation request to the server.
func (c *Client) CPUIntensive() (*benchpb.ComputeResponse, error) {
	return c.client.CPUIntensive(context.Background(), &benchpb.ComputeRequest{})
}

// MemoryIntensive sends a memory-intensive operation request to the server.
func (c *Client) MemoryIntensive() (*benchpb.MemoryResponse, error) {
	return c.client.MemoryIntensive(context.Background(), &benchpb.MemoryRequest{})
}

// DiskIO sends a disk I/O operation request to the server.
func (c *Client) DiskIO() (*benchpb.IOResponse, error) {
	return c.client.DiskIO(context.Background(), &benchpb.IORequest{})
}

// Mixed sends a mixed workload request to the server.
func (c *Client) Mixed() (*benchpb.MixedResponse, error) {
	return c.client.Mixed(context.Background(), &benchpb.MixedRequest{})
}

