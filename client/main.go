package main

import (
	"fmt"
	"log"
	"os"
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

	// Get benchmark parameters based on type
	switch benchmarkType {
	case HelloBenchmark:
		var numRequests int
		fmt.Print("Enter the number of sequential Hello requests: ")
		fmt.Scanf("%d", &numRequests)
		if numRequests <= 0 {
			numRequests = 1000 // default
		}
		
		start := time.Now()
		for i := 0; i < numRequests; i++ {
			if _, err := client.Hello(); err != nil {
				log.Fatalf("Failed to execute Hello benchmark: %v", err)
			}
		}
		duration := time.Since(start)
		log.Printf("Completed %d Hello requests in %v", numRequests, duration)
		log.Printf("Average time per request: %v", duration/time.Duration(numRequests))
		log.Printf("Requests per second: %.2f", float64(numRequests)/duration.Seconds())

	case CPUBenchmark:
		var iterations int32
		fmt.Print("Enter the number of iterations for CPU benchmark (default 1,000,000,000): ")
		fmt.Scanf("%d", &iterations)
		if iterations <= 0 {
			iterations = 1_000_000_000
		}

		start := time.Now()
		if _, err := client.CPUIntensive(iterations); err != nil {
			log.Fatalf("Failed to execute CPU benchmark: %v", err)
		}
		duration := time.Since(start)
		log.Printf("Completed CPU benchmark with %d iterations in %v", iterations, duration)

	case MemoryBenchmark:
		var sizeMB int32
		fmt.Print("Enter the memory size in MB (default 1024): ")
		fmt.Scanf("%d", &sizeMB)
		if sizeMB <= 0 {
			sizeMB = 1024 // 1GB default
		}

		start := time.Now()
		resp, err := client.MemoryIntensive(sizeMB)
		if err != nil {
			log.Fatalf("Failed to execute Memory benchmark: %v", err)
		}
		totalDuration := time.Since(start)
		
		log.Printf("Memory benchmark results:")
		log.Printf("Total size: %d MB", sizeMB)
		log.Printf("Pages accessed: %d", resp.PagesAccessed)
		log.Printf("Random access time: %v", time.Duration(resp.AccessTimeNs)*time.Nanosecond)
		if resp.PagesAccessed > 0 {
			log.Printf("Average time per page: %v", 
				time.Duration(resp.AccessTimeNs)*time.Nanosecond/time.Duration(resp.PagesAccessed))
		}
		log.Printf("Total execution time: %v", totalDuration)

	case DiskIOBenchmark:
		var fileSizeMB, numFiles int32
		fmt.Print("Enter the file size in MB (default 100): ")
		fmt.Scanf("%d", &fileSizeMB)
		if fileSizeMB <= 0 {
			fileSizeMB = 100
		}
		
		fmt.Print("Enter the number of files (default 5): ")
		fmt.Scanf("%d", &numFiles)
		if numFiles <= 0 {
			numFiles = 5
		}

		start := time.Now()
		if _, err := client.DiskIO(fileSizeMB, numFiles); err != nil {
			log.Fatalf("Failed to execute I/O benchmark: %v", err)
		}
		duration := time.Since(start)
		log.Printf("Completed I/O benchmark with %d files of %dMB each in %v", numFiles, fileSizeMB, duration)
		log.Printf("Total data processed: %dMB", fileSizeMB*numFiles)
		log.Printf("Average throughput: %.2f MB/s", float64(fileSizeMB*numFiles)/duration.Seconds())

	case MixedBenchmark:
		start := time.Now()
		if _, err := client.Mixed(); err != nil {
			log.Fatalf("Failed to execute Mixed benchmark: %v", err)
		}
		duration := time.Since(start)
		log.Printf("Completed Mixed benchmark in %v", duration)
	}
}