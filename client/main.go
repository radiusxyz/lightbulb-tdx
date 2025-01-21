package main

import (
	"log"
	"time"

	"github.com/radiusxyz/lightbulb-tdx/auction"
)

func main() {
	// Create a new client instance.
	client, err := auction.NewClient("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Add a new auction.
	log.Printf("%v", time.Now())
	client.AddAuction(1, "auction-1", time.Now().Add(1*time.Second), time.Now().Add(5*time.Second), "parameters")
}