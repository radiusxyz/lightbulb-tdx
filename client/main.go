package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	benchmark "github.com/radiusxyz/lightbulb-tdx/benchmark"
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

// Default benchmark parameters
const (
	DefaultHelloRequests   = 1000
	DefaultCPUIterations   = 1_000_000_000
	DefaultMemorySizeMB    = 1024 // 1GB
	DefaultIOFileSizeMB    = 100
	DefaultIONumFiles      = 5
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Initialize client
	client, err := benchmark.NewClient(os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	// Get benchmark type from user
	benchmarkType := getBenchmarkTypeFromUser()
	
	// Execute selected benchmark
	switch benchmarkType {
	case HelloBenchmark:
		runHelloBenchmark(client)
	case CPUBenchmark:
		runCPUBenchmark(client)
	case MemoryBenchmark:
		runMemoryBenchmark(client)
	case DiskIOBenchmark:
		runDiskIOBenchmark(client)
	case MixedBenchmark:
		runMixedBenchmark(client)
	}
}

// getBenchmarkTypeFromUser displays options and gets user's benchmark selection
func getBenchmarkTypeFromUser() BenchmarkType {
	fmt.Println("\nAvailable benchmark types:")
	fmt.Println("1. Hello (Network)")
	fmt.Println("2. CPU Intensive")
	fmt.Println("3. Memory Intensive")
	fmt.Println("4. Disk I/O")
	fmt.Println("5. Mixed")

	var choice int
	fmt.Print("\nSelect benchmark type (1-5): ")
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > 5 {
		log.Fatalf("Invalid benchmark type selection")
	}
	
	return BenchmarkType(choice - 1)
}

// getIntInput prompts for an integer with a default value
func getIntInput(prompt string, defaultValue int) int {
	var value int
	fmt.Printf("%s (default %d): ", prompt, defaultValue)
	_, err := fmt.Scanf("%d", &value)
	if err != nil || value <= 0 {
		return defaultValue
	}
	return value
}

// runHelloBenchmark executes the Hello benchmark
func runHelloBenchmark(client *benchmark.Client) {
	// Get parameters
	numRequests := getIntInput("Enter the number of sequential Hello requests", DefaultHelloRequests)
	
	// Execute benchmark
	start := time.Now()
	for range make([]struct{}, numRequests) {
		if _, err := client.Hello(); err != nil {
			log.Fatalf("Failed to execute Hello benchmark: %v", err)
		}
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("Hello benchmark results:")
	log.Printf("Completed %d requests in %v", numRequests, duration)
	log.Printf("Average time per request: %v", duration/time.Duration(numRequests))
	log.Printf("Requests per second: %.2f", float64(numRequests)/duration.Seconds())
}

// runCPUBenchmark executes the CPU intensive benchmark
func runCPUBenchmark(client *benchmark.Client) {
	// Get parameters
	iterations := int32(getIntInput("Enter the number of iterations for CPU benchmark", DefaultCPUIterations))
	
	// Execute benchmark
	start := time.Now()
	if _, err := client.CPUIntensive(iterations); err != nil {
		log.Fatalf("Failed to execute CPU benchmark: %v", err)
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("CPU benchmark results:")
	log.Printf("Completed %d iterations in %v", iterations, duration)
	log.Printf("Iterations per second: %.2f", float64(iterations)/duration.Seconds())
}

// runMemoryBenchmark executes the Memory intensive benchmark
func runMemoryBenchmark(client *benchmark.Client) {
	// Get parameters
	sizeMB := int32(getIntInput("Enter the memory size in MB", DefaultMemorySizeMB))
	
	// Execute benchmark
	start := time.Now()
	resp, err := client.MemoryIntensive(sizeMB)
	if err != nil {
		log.Fatalf("Failed to execute Memory benchmark: %v", err)
	}
	totalDuration := time.Since(start)
	
	// Report results
	log.Printf("Memory benchmark results:")
	log.Printf("Total size: %d MB", sizeMB)
	log.Printf("Pages accessed: %d", resp.PagesAccessed)
	log.Printf("Random access time: %v", time.Duration(resp.AccessTimeNs)*time.Nanosecond)
	if resp.PagesAccessed > 0 {
		log.Printf("Average time per page: %v", 
			time.Duration(resp.AccessTimeNs)*time.Nanosecond/time.Duration(resp.PagesAccessed))
	}
	log.Printf("Total execution time: %v", totalDuration)
}

// runDiskIOBenchmark executes the Disk I/O benchmark
func runDiskIOBenchmark(client *benchmark.Client) {
	// Get parameters
	fileSizeMB := int32(getIntInput("Enter the file size in MB", DefaultIOFileSizeMB))
	numFiles := int32(getIntInput("Enter the number of files", DefaultIONumFiles))
	
	// Execute benchmark
	start := time.Now()
	if _, err := client.DiskIO(fileSizeMB, numFiles); err != nil {
		log.Fatalf("Failed to execute I/O benchmark: %v", err)
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("Disk I/O benchmark results:")
	log.Printf("Completed with %d files of %d MB each in %v", numFiles, fileSizeMB, duration)
	log.Printf("Total data processed: %d MB", fileSizeMB*numFiles)
	log.Printf("Average throughput: %.2f MB/s", float64(fileSizeMB*numFiles)/duration.Seconds())
}

// runMixedBenchmark executes the Mixed benchmark
func runMixedBenchmark(client *benchmark.Client) {
	// Execute benchmark
	start := time.Now()
	if _, err := client.Mixed(); err != nil {
		log.Fatalf("Failed to execute Mixed benchmark: %v", err)
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("Mixed benchmark results:")
	log.Printf("Completed in %v", duration)
}