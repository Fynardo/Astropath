package claude

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
)

// Package claude provides a way for interacting with Claude Code.


type PromptParams struct {
}

type ReviewerParams struct {
	BranchName string // Name of the pull request to review
}


// streamReader reads from a pipe and processes stream-json output
func streamReader(reader io.Reader, done chan<- error) {
	defer close(done)
	
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		
		// Try to parse as JSON to extract content
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err == nil {
			// Extract and print content if it exists
			if msgType, ok := jsonData["type"].(string); ok && msgType == "user" {
				fmt.Println(jsonData["message"])
			}
			if content, ok := jsonData["content"].(string); ok && content != "" {
				// fmt.Print(content)
			}
		} else {
			// If not JSON, just print the line as-is
			fmt.Println(line)
		}
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

// RunAgentWithStreaming spawns a Claude agent with streaming output.
// It returns a channel that will receive an error when the agent finishes.
func RunAgentWithStreaming(prompt string) <-chan error {
	done := make(chan error, 1)
	
	go func() {
		defer close(done)
		
		fmt.Println("Starting Claude agent with streaming...")
		fmt.Println("Using prompt:\n=====\n", prompt)
		fmt.Println("=====")

		// Create command with pipes
		cmd := exec.Command("claude", "--verbose", "-p", prompt, "--output-format", "stream-json")
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
		// Create pipe for stdout
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			done <- fmt.Errorf("failed to create stdout pipe: %v", err)
			return
		}
		
		// Start the command
		if err := cmd.Start(); err != nil {
			done <- fmt.Errorf("failed to start Claude agent: %v", err)
			return
		}
		
		// Start stream reader goroutine
		streamDone := make(chan error, 1)
		go streamReader(stdout, streamDone)
		
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
		
		done <- finalErr
	}()
	
	return done
}
