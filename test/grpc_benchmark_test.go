package test

import (
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/radiusxyz/lightbulb-tdx/benchmark"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load("../.env"); err != nil {
		panic("Error loading .env file: " + err.Error())
	}
}

// BenchmarkGRPCHello implements different concurrency patterns for benchmarking
func BenchmarkGRPCHello(b *testing.B) {
	// Test cases with different worker counts
	workerCounts := []int{1, 10, 50, 100, 200}
	requestCount := 100000

	for _, workers := range workerCounts {
		b.Run(fmt.Sprintf("Workers=%d,Requests=%d", workers, requestCount), func(b *testing.B) {
			client, err := benchmark.NewClient(os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("PORT"))
			if err != nil {
				b.Fatalf("Failed to create client: %v", err)
			}
			defer client.Close()

			start := time.Now() // Record start time
			b.ResetTimer()

			var wg sync.WaitGroup
			taskChan := make(chan int, b.N)

			// Launch workers
			for i := 0; i < workers; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					for range taskChan {
						_, err := client.Hello()
						if err != nil {
							b.Errorf("Failed to call Hello: %v", err)
							return
						}
					}
				}()
			}

			// Send tasks to workers
			for i := 0; i < b.N; i++ {
				taskChan <- i
			}
			close(taskChan)

			wg.Wait()
			duration := time.Since(start)
			b.ReportMetric(float64(duration.Milliseconds())/float64(b.N), "ms/op")
			b.ReportMetric(float64(b.N)/duration.Seconds(), "requests/sec")
		})
	}
}

// TestGRPCHello is a basic functionality test
func TestGRPCHello(t *testing.T) {
	client, err := benchmark.NewClient(os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("PORT"))
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	_, err = client.Hello()
	if err != nil {
		t.Errorf("Hello call failed: %v", err)
	}
}

// BenchmarkGRPCHelloSequential benchmarks sequential requests
func BenchmarkGRPCHelloSequential(b *testing.B) {
	client, err := benchmark.NewClient(os.Getenv("SERVER_ADDRESS") + ":" + os.Getenv("PORT"))
	if err != nil {
		b.Fatalf("Failed to create client: %v", err)
	}
	defer client.Close()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := client.Hello()
		if err != nil {
			b.Errorf("Failed to call Hello: %v", err)
			return
		}
	}
} 