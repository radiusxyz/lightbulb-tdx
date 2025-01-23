package main

import (
	"log"
	"net"
	"os"
	"runtime"
	"runtime/pprof"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/radiusxyz/lightbulb-tdx/auction"
	"github.com/radiusxyz/lightbulb-tdx/tdx"

	attestpb "github.com/radiusxyz/lightbulb-tdx/proto/attest"
	auctionpb "github.com/radiusxyz/lightbulb-tdx/proto/auction"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	
	// Enable profiling if PROFILING is set to true
	if os.Getenv("PROFILING") == "true" {
		// Start CPU profiling
		cpuFile, err := os.Create("cpu.prof")
		if err != nil {
			log.Fatalf("could not create CPU profile: %v", err)
		}
		defer cpuFile.Close()

		if err := pprof.StartCPUProfile(cpuFile); err != nil {
			log.Fatalf("could not start CPU profile: %v", err)
		}
		defer pprof.StopCPUProfile()

		// Start memory profiling
		memFile, err := os.Create("mem.prof")
		if err != nil {
			log.Fatalf("could not create memory profile: %v", err)
		}
		defer memFile.Close()

		runtime.GC() // get up-to-date statistics

		if err := pprof.WriteHeapProfile(memFile); err != nil {
			log.Fatalf("could not write memory profile: %v", err)
		}
		defer pprof.StopCPUProfile()
	}

	// Listen on the specified port
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

	// Start the server
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
