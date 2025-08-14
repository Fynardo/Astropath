package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

// Package claude provides a way for interacting with Claude Code.

// StreamBuffer holds captured streaming output with thread-safe access
type StreamBuffer struct {
	mutex     sync.RWMutex
	buffer    []byte
	maxSize   int
	jsonLines []string
}

// NewStreamBuffer creates a new stream buffer with specified maximum size
func NewStreamBuffer(maxSize int) *StreamBuffer {
	if maxSize <= 0 {
		maxSize = 1024 * 1024 // Default 1MB
	}
	return &StreamBuffer{
		buffer:    make([]byte, 0),
		maxSize:   maxSize,
		jsonLines: make([]string, 0),
	}
}

// Write implements io.Writer interface for thread-safe writing
func (sb *StreamBuffer) Write(data []byte) (int, error) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	
	// Append to buffer with size limit
	if len(sb.buffer)+len(data) > sb.maxSize {
		// Remove old data from beginning to make room
		excess := len(sb.buffer) + len(data) - sb.maxSize
		if excess >= len(sb.buffer) {
			sb.buffer = sb.buffer[:0]
		} else {
			sb.buffer = sb.buffer[excess:]
		}
	}
	
	sb.buffer = append(sb.buffer, data...)
	return len(data), nil
}

// GetBuffer returns a copy of the current buffer content
func (sb *StreamBuffer) GetBuffer() []byte {
	sb.mutex.RLock()
	defer sb.mutex.RUnlock()
	
	result := make([]byte, len(sb.buffer))
	copy(result, sb.buffer)
	return result
}

// AddJSONLine stores a parsed JSON line from the stream
func (sb *StreamBuffer) AddJSONLine(line string) {
	sb.mutex.Lock()
	defer sb.mutex.Unlock()
	
	sb.jsonLines = append(sb.jsonLines, line)
	// Keep only last 1000 lines to prevent memory overflow
	if len(sb.jsonLines) > 1000 {
		sb.jsonLines = sb.jsonLines[len(sb.jsonLines)-1000:]
	}
}

// GetJSONLines returns a copy of parsed JSON lines
func (sb *StreamBuffer) GetJSONLines() []string {
	sb.mutex.RLock()
	defer sb.mutex.RUnlock()
	
	result := make([]string, len(sb.jsonLines))
	copy(result, sb.jsonLines)
	return result
}

// StreamResult contains both error status and buffer access for enhanced streaming
type StreamResult struct {
	Error  error
	Buffer *StreamBuffer
}

type PromptParams struct {
}

type ReviewerParams struct {
	BranchName string // Name of the pull request to review
}

// TeeWriter writes to multiple writers simultaneously
type TeeWriter struct {
	writers []io.Writer
}

func NewTeeWriter(writers ...io.Writer) *TeeWriter {
	return &TeeWriter{writers: writers}
}

func (tw *TeeWriter) Write(data []byte) (int, error) {
	for _, w := range tw.writers {
		if _, err := w.Write(data); err != nil {
			return 0, err
		}
	}
	return len(data), nil
}

// streamReader reads from a pipe and processes stream-json output
func streamReader(reader io.Reader, buffer *StreamBuffer, done chan<- error) {
	defer close(done)
	
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Try to parse as JSON to extract content
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err == nil {
			buffer.AddJSONLine(line)
			
			// Extract and print content if it exists
			if content, ok := jsonData["content"].(string); ok && content != "" {
				fmt.Print(content)
			}
		} else {
			// If not JSON, just print the line as-is
			fmt.Println(line)
		}
		
		// Also write raw line to buffer
		buffer.Write([]byte(line + "\n"))
	}
	
	if err := scanner.Err(); err != nil {
		done <- fmt.Errorf("error reading stream: %v", err)
		return
	}
	
	done <- nil
}

// RunAgent spawns a Claude agent using the 'claude -p' command with the specified prompt type.
// It returns a channel that will receive a signal when the agent finishes.
func RunAgent(prompt string) <-chan error {
	done := make(chan error, 1)
	
	go func() {
		fmt.Println("Starting Claude agent...")
		fmt.Println("Using prompt:\n=====\n", prompt)
		fmt.Println("=====")

		cmd := exec.Command("claude", "--verbose", "-p", prompt, "--output-format", "stream-json")
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		
		err := cmd.Run()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error running Claude agent: %v\n", err)
		} else {
			fmt.Println("Claude agent finished.")
		}
		
		done <- err
		close(done)
	}()
	
	return done
}

// RunAgentWithStreaming spawns a Claude agent with streaming output capture.
// It returns a channel that will receive a StreamResult when the agent finishes.
func RunAgentWithStreaming(prompt string) <-chan StreamResult {
	done := make(chan StreamResult, 1)
	
	go func() {
		defer close(done)
		
		fmt.Println("Starting Claude agent with streaming...")
		fmt.Println("Using prompt:\n=====\n", prompt)
		fmt.Println("=====")

		// Create stream buffer
		buffer := NewStreamBuffer(1024 * 1024) // 1MB buffer
		
		// Create command with pipes
		cmd := exec.Command("claude", "--verbose", "-p", prompt, "--output-format", "stream-json")
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		// Create pipe for stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			done <- StreamResult{Error: fmt.Errorf("failed to create stdout pipe: %v", err), Buffer: buffer}
			return
		}
		
		// Start the command
		if err := cmd.Start(); err != nil {
			done <- StreamResult{Error: fmt.Errorf("failed to start Claude agent: %v", err), Buffer: buffer}
			return
		}
		
		// Start stream reader goroutine
		streamDone := make(chan error, 1)
		go streamReader(stdout, buffer, streamDone)
		
		// Wait for both command and stream reader to finish
		cmdErr := cmd.Wait()
		streamErr := <-streamDone
		
		// Determine final error
		var finalErr error
		if cmdErr != nil {
			finalErr = fmt.Errorf("Claude agent error: %v", cmdErr)
		} else if streamErr != nil {
			finalErr = streamErr
		}
		
		if finalErr != nil {
			fmt.Fprintf(os.Stderr, "Error running Claude agent: %v\n", finalErr)
		} else {
			fmt.Println("Claude agent finished.")
		}
		
		done <- StreamResult{Error: finalErr, Buffer: buffer}
	}()
	
	return done
}
