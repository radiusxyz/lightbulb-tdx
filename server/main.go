package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/radiusxyz/lightbulb-tdx/auction"
	"github.com/radiusxyz/lightbulb-tdx/benchmark"
	"github.com/radiusxyz/lightbulb-tdx/tdx"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
	benchmarkpb "github.com/radiusxyz/lightbulb-tdx/proto/benchmark"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Listen on the specified port
	lis, err := net.Listen("tcp", ":"+os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on port %s in %s environment", lis.Addr(), os.Getenv("ENV"))

	// Create gRPC server and TDX client
	grpcServer := grpc.NewServer()
	tdxClient := tdx.NewTDXClient()

	// Create and register services
	attestServer := tdx.NewServer(tdxClient)
	auctionServer := auction.NewServer()
	benchmarkServer, err := benchmark.NewServer()
	if err != nil {
		log.Fatalf("Failed to create benchmark server: %v", err)
	}
	attestpb.RegisterAttestServiceServer(grpcServer, attestServer)
	auctionpb.RegisterAuctionServiceServer(grpcServer, auctionServer)
	benchmarkpb.RegisterBenchmarkServiceServer(grpcServer, benchmarkServer)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	// Graceful shutdown setup
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Run the server in a goroutine
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-stop
	log.Println("Shutting down server gracefully...")

	// Create a context with timeout for shutdown
	_, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Gracefully stop the gRPC server
	grpcServer.GracefulStop()

	// Perform any additional cleanup tasks if necessary
	log.Println("Server stopped")
}
