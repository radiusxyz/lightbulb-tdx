package tdx

import (
	"bufio"
	"crypto"
	"fmt"
	"os"
	"sync"
)

// RtmrProvider is a provider for RTMR values.
type RtmrProvider struct {
	rtmrs			  [4][]byte   // RTMR values
	mu                sync.Mutex  // Protects lastProcessedLine
	lastProcessedLine int         // The last processed line number
	logPath           string      // Path to the IMA log file
	hashAlgo          crypto.Hash // Hash algorithm to use
}

// NewRtmrProvider creates a new RtmrProvider.
func NewRtmrProvider(logPath string, hashAlgo crypto.Hash) *RtmrProvider {
	return &RtmrProvider{
		lastProcessedLine: 0,
		logPath:           logPath,
		hashAlgo:          hashAlgo,
	}
}

// DefaultRtmrProvider creates a new IMARtmrProvider with default settings.
func DefaultRtmrProvider() *RtmrProvider {
	logPath := os.Getenv("IMA_LOG_PATH")
	return NewRtmrProvider(logPath, crypto.SHA384)
}

// UpdateImaRtmr reads the IMA log file and updates the RTMR values.
func (e *RtmrProvider) UpdateImaRtmr() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Open the IMA log file
	file, err := os.Open(e.logPath)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer file.Close()

	// Read the file line by line
	scanner := bufio.NewScanner(file)
	currentLine := 0
	rtmrIndex := 2
	for scanner.Scan() {
		currentLine++

		// Skip lines that have already been processed
		if currentLine <= e.lastProcessedLine {
			continue
		}

		// Process the current line
		eventLog := scanner.Bytes()
		
		// Compute the hash of the event log
		hasher := e.hashAlgo.New()
		hasher.Write(eventLog)
		digest := hasher.Sum(nil)

		// Extend the RTMR with the new digest
		err := e.ExtendRtmr(rtmrIndex, digest)
		if err != nil {
			return fmt.Errorf("failed to extend event log on line %d: %w", currentLine, err)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read log file: %w", err)
	}

	// Update the last processed line
	e.lastProcessedLine = currentLine
	return nil
}

// GetLastProcessedLine retrieves the last processed line number.
func (e *RtmrProvider) GetLastProcessedLine() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.lastProcessedLine
}

// SetLastProcessedLine sets the last processed line number manually.
func (e *RtmrProvider) SetLastProcessedLine(line int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastProcessedLine = line
}

func (e *RtmrProvider) GetRtmrValues() [4][]byte {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.rtmrs
}

// ExtendRtmr concatenates the current RTMR value with the new digest and updates the RTMR with the hash of the result.
func (e *RtmrProvider) ExtendRtmr(index int, digest []byte) error {
	// Validate RTMR index
	if index < 0 || index >= len(e.rtmrs) {
		return fmt.Errorf("invalid index: %d", index)
	}

	// Check if hash algorithm is available
	if !e.hashAlgo.Available() {
		return fmt.Errorf("hash algorithm %v is not available", e.hashAlgo)
	}

	// Create a new hash instance
	hasher := e.hashAlgo.New()

	// Concatenate current RTMR value with the new digest
	currentRTMR := e.rtmrs[index]
	hasher.Write(currentRTMR) // Add the current RTMR
	hasher.Write(digest)      // Add the new digest

	// Compute the new hash
	newRTMR := hasher.Sum(nil)

	// Update the RTMR with the new hash
	e.rtmrs[index] = newRTMR

	return nil
}