package attest

import (
	"bufio"
	"crypto"
	"fmt"
	"os"
	"sync"

	"github.com/google/go-tdx-guest/rtmr"
)

// IMAEventLogExtender manages the state for extending IMA event logs.
type IMAEventLogExtender struct {
	mu                sync.Mutex  // Protects lastProcessedLine
	lastProcessedLine int         // The last processed line number
	logPath           string      // Path to the IMA log file
	hashAlgo          crypto.Hash // Hash algorithm to use
}

// NewIMAEventLogExtender creates a new IMAEventLogExtender.
func NewIMAEventLogExtender(logPath string, hashAlgo crypto.Hash) *IMAEventLogExtender {
	return &IMAEventLogExtender{
		lastProcessedLine: 0,
		logPath:           logPath,
		hashAlgo:          hashAlgo,
	}
}

// DefaultIMAEventLogExtender creates a new IMAEventLogExtender with default settings.
func DefaultIMAEventLogExtender() *IMAEventLogExtender {
	logPath := "/sys/kernel/security/ima/ascii_runtime_measurements"
	return NewIMAEventLogExtender(logPath, crypto.SHA384)
}

// ExtendLogs reads the IMA event log and extends only new entries.
func (e *IMAEventLogExtender) ExtendLogs(rtmrIndex int) error {
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
	for scanner.Scan() {
		currentLine++

		// Skip lines that have already been processed
		if currentLine <= e.lastProcessedLine {
			continue
		}

		// Process the current line
		eventLog := scanner.Bytes()
		if err := rtmr.ExtendEventLog(rtmrIndex, e.hashAlgo, eventLog); err != nil {
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
func (e *IMAEventLogExtender) GetLastProcessedLine() int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.lastProcessedLine
}

// SetLastProcessedLine sets the last processed line number manually.
func (e *IMAEventLogExtender) SetLastProcessedLine(line int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.lastProcessedLine = line
}