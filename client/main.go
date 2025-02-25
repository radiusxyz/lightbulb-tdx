package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
	benchmark "github.com/radiusxyz/lightbulb-tdx/benchmark"
)

// BenchmarkType represents different types of benchmarks
type BenchmarkType int

const (
	SequentialNetworkBenchmark BenchmarkType = iota
	ConcurrentNetworkBenchmark
	CPUBenchmark
	MemoryBenchmark
	DiskIOBenchmark
	MixedBenchmark
)

// Default benchmark parameters
const (
	DefaultSequentialRequests   = 100_000
	DefaultConcurrentClients    = 100
	DefaultRequestsPerClient    = 10_000
	DefaultCPUIterations        = 1_000_000_000
	DefaultMemorySizeMB         = 1024 // 1GB
	DefaultIOFileSizeMB         = 100
	DefaultIONumFiles           = 5
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
	case SequentialNetworkBenchmark:
		benchmarkSequentialNetwork(client)
	case ConcurrentNetworkBenchmark:
		benchmarkConcurrentNetwork(client)
	case CPUBenchmark:
		benchmarkCPU(client)
	case MemoryBenchmark:
		benchmarkMemory(client)
	case DiskIOBenchmark:
		benchmarkDiskIO(client)
	case MixedBenchmark:
		benchmarkMixed(client)
	}
}

// getBenchmarkTypeFromUser displays options and gets user's benchmark selection
func getBenchmarkTypeFromUser() BenchmarkType {
	fmt.Println("\nAvailable benchmark types:")
	fmt.Println("1. Network - Sequential Requests")
	fmt.Println("2. Network - Concurrent Requests")
	fmt.Println("3. CPU Intensive")
	fmt.Println("4. Memory Intensive")
	fmt.Println("5. Disk I/O")
	fmt.Println("6. Mixed (All Resources)")

	var choice int
	fmt.Print("\nSelect benchmark type (1-6): ")
	_, err := fmt.Scanf("%d", &choice)
	if err != nil || choice < 1 || choice > 6 {
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

// benchmarkSequentialNetwork tests network performance with sequential requests
func benchmarkSequentialNetwork(client *benchmark.Client) {
	// Get parameters
	numRequests := getIntInput("Enter the number of sequential requests", DefaultSequentialRequests)
	
	// Execute benchmark
	start := time.Now()
	for range numRequests {
		if _, err := client.Hello(); err != nil {
			log.Fatalf("Failed to execute network benchmark: %v", err)
		}
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("\nSequential Network Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Completed %d requests in %v", numRequests, duration)
	log.Printf("Average time per request: %v", duration/time.Duration(numRequests))
	log.Printf("Throughput: %.2f requests/second", float64(numRequests)/duration.Seconds())
}

// benchmarkConcurrentNetwork tests network performance with concurrent requests
func benchmarkConcurrentNetwork(client *benchmark.Client) {
	// Get parameters
	numClients := getIntInput("Enter the number of concurrent clients", DefaultConcurrentClients)
	requestsPerClient := getIntInput("Enter the number of requests per client", DefaultRequestsPerClient)
	totalRequests := numClients * requestsPerClient
	
	// Channel to collect results from each goroutine
	resultCh := make(chan time.Duration, numClients)
	errorCh := make(chan error, totalRequests) // Channel to collect errors
	
	// Create a wait group to wait for all goroutines to complete
	var wg sync.WaitGroup
	wg.Add(numClients)
	
	// Start timer
	start := time.Now()
	
	// Launch client goroutines
	for i := range numClients {
		go func(clientID int) {
			defer wg.Done()
			
			clientStart := time.Now()
			for j := range requestsPerClient {
				if _, err := client.Hello(); err != nil {
					errorCh <- fmt.Errorf("client %d: failed to execute request %d: %v", clientID, j, err)
					continue
				}
			}
			clientDuration := time.Since(clientStart)
			resultCh <- clientDuration
		}(i)
	}
	
	// Close the result channel once all goroutines are done
	go func() {
		wg.Wait()
		close(resultCh)
		close(errorCh)
	}()
	
	// Collect results
	var totalClientDuration time.Duration
	minDuration := time.Hour
	maxDuration := time.Duration(0)
	clientDurations := make([]time.Duration, 0, numClients)
	
	for duration := range resultCh {
		clientDurations = append(clientDurations, duration)
		totalClientDuration += duration
		if duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
	}
	
	// Process any errors
	errorCount := 0
	for err := range errorCh {
		errorCount++
		log.Printf("Error: %v", err)
	}
	
	// Calculate total elapsed time
	totalDuration := time.Since(start)
	successfulRequests := totalRequests - errorCount
	
	// Report results
	log.Printf("\nConcurrent Network Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Total clients: %d", numClients)
	log.Printf("Requests per client: %d", requestsPerClient)
	log.Printf("Total requests: %d", totalRequests)
	log.Printf("Successful requests: %d (%.2f%%)", successfulRequests, float64(successfulRequests)*100/float64(totalRequests))
	log.Printf("Failed requests: %d", errorCount)
	
	if len(clientDurations) > 0 {
		log.Printf("Average client duration: %v", totalClientDuration/time.Duration(len(clientDurations)))
		log.Printf("Fastest client: %v", minDuration)
		log.Printf("Slowest client: %v", maxDuration)
	}
	
	log.Printf("Total elapsed time: %v", totalDuration)
	log.Printf("Overall throughput: %.2f requests/second", float64(successfulRequests)/totalDuration.Seconds())
}

// benchmarkCPU tests CPU performance
func benchmarkCPU(client *benchmark.Client) {
	// Get parameters
	iterations := int32(getIntInput("Enter the number of iterations for CPU benchmark", DefaultCPUIterations))
	
	// Execute benchmark
	start := time.Now()
	if _, err := client.CPUIntensive(iterations); err != nil {
		log.Fatalf("Failed to execute CPU benchmark: %v", err)
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("\nCPU Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Completed %d iterations in %v", iterations, duration)
	log.Printf("Iterations per second: %.2f", float64(iterations)/duration.Seconds())
}

// benchmarkMemory tests memory performance
func benchmarkMemory(client *benchmark.Client) {
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
	log.Printf("\nMemory Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Total size: %d MB", sizeMB)
	log.Printf("Pages accessed: %d", resp.PagesAccessed)
	log.Printf("Random access time: %v", time.Duration(resp.AccessTimeNs)*time.Nanosecond)
	if resp.PagesAccessed > 0 {
		log.Printf("Average time per page: %v", 
			time.Duration(resp.AccessTimeNs)*time.Nanosecond/time.Duration(resp.PagesAccessed))
	}
	log.Printf("Total execution time: %v", totalDuration)
}

// benchmarkDiskIO tests disk I/O performance
func benchmarkDiskIO(client *benchmark.Client) {
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
	log.Printf("\nDisk I/O Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Completed with %d files of %d MB each in %v", numFiles, fileSizeMB, duration)
	log.Printf("Total data processed: %d MB", fileSizeMB*numFiles)
	log.Printf("Average throughput: %.2f MB/s", float64(fileSizeMB*numFiles)/duration.Seconds())
}

// benchmarkMixed tests mixed workload performance
func benchmarkMixed(client *benchmark.Client) {
	// Execute benchmark
	start := time.Now()
	if _, err := client.Mixed(); err != nil {
		log.Fatalf("Failed to execute Mixed benchmark: %v", err)
	}
	duration := time.Since(start)
	
	// Report results
	log.Printf("\nMixed Benchmark Results:")
	log.Printf("-----------------------------------")
	log.Printf("Completed in %v", duration)
}