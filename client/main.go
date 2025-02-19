package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	auction "github.com/radiusxyz/lightbulb-tdx/benchmark"
)

// BenchmarkType represents different types of benchmarks
type BenchmarkType int

const (
	HelloBenchmark BenchmarkType = iota
	CPUBenchmark
	MemoryBenchmark
	DiskIOBenchmark
	MixedBenchmark
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

	// Get the benchmark type from user
	fmt.Println("\nAvailable benchmark types:")
	fmt.Println("1. Hello (Network)")
	fmt.Println("2. CPU Intensive")
	fmt.Println("3. Memory Intensive")
	fmt.Println("4. Disk I/O")
	fmt.Println("5. Mixed")

	var benchmarkChoice int
	fmt.Print("\nSelect benchmark type (1-5): ")
	_, err = fmt.Scanf("%d", &benchmarkChoice)
	if err != nil || benchmarkChoice < 1 || benchmarkChoice > 5 {
		log.Fatalf("Invalid benchmark type selection")
	}
	benchmarkType := BenchmarkType(benchmarkChoice - 1)

	// Get the number of workers and iterations from user input
	var numWorkers, iterations int
	
	fmt.Print("Enter the number of concurrent workers: ")
	_, err = fmt.Scanf("%d", &numWorkers)
	if err != nil {
		log.Fatalf("Input error: %v", err)
	}
	if numWorkers < 1 {
		log.Fatalf("Number of workers must be at least 1")
	}

	fmt.Print("Enter the total number of requests: ")
	_, err = fmt.Scanf("%d", &iterations)
	if err != nil {
		log.Fatalf("Input error: %v", err)
	}
	if iterations < 1 {
		log.Fatalf("Number of requests must be at least 1")
	}

	log.Printf("Starting %s benchmark with %d workers, sending %d requests...",
		getBenchmarkName(benchmarkType), numWorkers, iterations)

	// Measure execution time
	start := time.Now()

	var wg sync.WaitGroup
	taskChan := make(chan int, iterations)

	// Launch workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range taskChan {
				var err error
				switch benchmarkType {
				case HelloBenchmark:
					_, err = client.Hello()
				case CPUBenchmark:
					_, err = client.CPUIntensive()
				case MemoryBenchmark:
					_, err = client.MemoryIntensive()
				case DiskIOBenchmark:
					_, err = client.DiskIO()
				case MixedBenchmark:
					_, err = client.Mixed()
				}
				
				if err != nil {
					log.Printf("Failed to execute benchmark: %v", err)
					return
				}
			}
		}()
	}

	// Send tasks to workers
	for i := 0; i < iterations; i++ {
		taskChan <- i
	}
	close(taskChan)

	// Wait for all workers to complete
	wg.Wait()
	
	duration := time.Since(start)
	log.Printf("Total execution time: %v", duration)
	log.Printf("Average execution time per request: %v", duration/time.Duration(iterations))
	log.Printf("Requests per second: %.2f", float64(iterations)/duration.Seconds())
}

func getBenchmarkName(bt BenchmarkType) string {
	switch bt {
	case HelloBenchmark:
		return "Hello"
	case CPUBenchmark:
		return "CPU Intensive"
	case MemoryBenchmark:
		return "Memory Intensive"
	case DiskIOBenchmark:
		return "Disk I/O"
	case MixedBenchmark:
		return "Mixed"
	default:
		return "Unknown"
	}
}