package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/radiusxyz/lightbulb-tdx/auction"
)

func main() {
	// Load environment variables from .env file.
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create a new client instance.
	client, err := auction.NewClient(os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Add a new auction.
	log.Printf("%v", time.Now())
	client.AddAuction(1, "auction-1", time.Now().Add(1*time.Second), time.Now().Add(5*time.Second), "parameters")
}