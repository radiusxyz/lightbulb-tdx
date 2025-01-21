package main

import (
	"log"
	"net"
	"os"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/radiusxyz/lightbulb-tdx/auction"
	"github.com/radiusxyz/lightbulb-tdx/tdx"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	lis, err := net.Listen("tcp", ":" + os.Getenv("PORT"))
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

	attestpb.RegisterAttestServiceServer(grpcServer, attestServer)
	auctionpb.RegisterAuctionServiceServer(grpcServer, auctionServer)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
