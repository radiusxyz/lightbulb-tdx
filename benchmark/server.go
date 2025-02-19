package benchmark

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	benchpb "github.com/radiusxyz/lightbulb-tdx/proto/benchmark"
)

type Server struct {
	benchpb.UnimplementedBenchmarkServiceServer
	tempDir string
}

// NewServer initializes a new gRPC server instance.
func NewServer() (*Server, error) {
	// Create temporary directory for I/O tests
	tempDir, err := os.MkdirTemp("", "benchmark-*")
	if err != nil {
		return nil, err
	}
	return &Server{tempDir: tempDir}, nil
}

// Cleanup removes temporary files
func (s *Server) Cleanup() error {
	return os.RemoveAll(s.tempDir)
}

// Hello tests basic network performance
func (s *Server) Hello(ctx context.Context, req *benchpb.HelloRequest) (*benchpb.HelloResponse, error) {
	return &benchpb.HelloResponse{Message: "Hello, World!"}, nil
}

// CPUIntensive tests CPU performance with heavy computation
func (s *Server) CPUIntensive(ctx context.Context, req *benchpb.ComputeRequest) (*benchpb.ComputeResponse, error) {
	result := 0.0
	// Perform CPU-intensive calculations
	for i := 0; i < 1000000; i++ {
		result += math.Sqrt(float64(i)) * math.Sin(float64(i))
	}
	return &benchpb.ComputeResponse{Result: result}, nil
}

// MemoryIntensive tests memory allocation and access patterns
func (s *Server) MemoryIntensive(ctx context.Context, req *benchpb.MemoryRequest) (*benchpb.MemoryResponse, error) {
	// Allocate and manipulate large chunks of memory
	size := 50 * 1024 * 1024 // 50MB
	data := make([]byte, size)
	
	// Perform memory operations
	for i := 0; i < size; i += 4096 {
		data[i] = byte(i)
	}

	// Calculate hash to prevent optimization
	hash := sha256.Sum256(data)
	return &benchpb.MemoryResponse{Hash: hash[:]}, nil
}

// DiskIO tests I/O performance
func (s *Server) DiskIO(ctx context.Context, req *benchpb.IORequest) (*benchpb.IOResponse, error) {
	// Create multiple files and perform concurrent I/O operations
	const (
		fileSize    = 10 * 1024 * 1024 // 10MB
		numFiles    = 5
		bufferSize  = 64 * 1024 // 64KB chunks
	)

	var wg sync.WaitGroup
	errChan := make(chan error, numFiles)

	for i := 0; i < numFiles; i++ {
		wg.Add(1)
		go func(fileNum int) {
			defer wg.Done()

			// Generate data
			data := make([]byte, fileSize)
			for j := range data {
				data[j] = byte(j % 256)
			}

			// Write file
			filename := filepath.Join(s.tempDir, fmt.Sprintf("test-%d.dat", fileNum))
			if err := os.WriteFile(filename, data, 0666); err != nil {
				errChan <- err
				return
			}

			// Read and verify
			readData, err := os.ReadFile(filename)
			if err != nil {
				errChan <- err
				return
			}

			// Calculate hash
			sha256.Sum256(readData)
			
			// Cleanup
			os.Remove(filename)
		}(i)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return nil, err
		}
	}

	return &benchpb.IOResponse{Success: true}, nil
}

// Mixed combines different types of workloads
func (s *Server) Mixed(ctx context.Context, req *benchpb.MixedRequest) (*benchpb.MixedResponse, error) {
	var wg sync.WaitGroup
	numCPUs := runtime.NumCPU()
	
	// CPU workload
	result := 0.0
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100000; j++ {
				result += math.Sqrt(float64(j))
			}
		}()
	}

	// Memory workload
	data := make([][]byte, numCPUs)
	for i := 0; i < numCPUs; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			data[i] = make([]byte, 1024*1024) // 1MB per CPU
			for j := range data[i] {
				data[i][j] = byte(j)
			}
		}(i)
	}

	// I/O workload
	wg.Add(1)
	go func() {
		defer wg.Done()
		filename := filepath.Join(s.tempDir, "mixed-test.dat")
		data := make([]byte, 1024*1024) // 1MB
		os.WriteFile(filename, data, 0666)
		os.ReadFile(filename)
		os.Remove(filename)
	}()

	wg.Wait()
	return &benchpb.MixedResponse{Success: true}, nil
}
