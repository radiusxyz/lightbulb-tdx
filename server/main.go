package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/radiusxyz/lightbulb-tdx/server/attest"
	"github.com/radiusxyz/lightbulb-tdx/server/auction"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Println("Server listening on port 50051...")

	// Create gRPC server and TDX client
	grpcServer := grpc.NewServer()
	tdxClient := attest.NewTDXClientWrapper()

	// Create and register services
	attestServer := attest.NewServer(tdxClient)
	auctionServer := auction.NewServer()

	attestpb.RegisterAttestServiceServer(grpcServer, attestServer)
	auctionpb.RegisterAuctionServiceServer(grpcServer, auctionServer)

	// Enable reflection for debugging
	reflection.Register(grpcServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
