package claude

import (
	"fmt"
	"os"
	"os/exec"
)

// Package claude provides a way for interacting with Claude Code.
type PromptParams struct {
}

type ReviewerParams struct {
	BranchName string // Name of the pull request to review
}

// RunAgent spawns a Claude agent using the 'claude -p' command with the specified prompt type.
// It returns a channel that will receive a signal when the agent finishes.
func RunAgent(prompt string) <-chan error {
	done := make(chan error, 1)
	
	go func() {
		fmt.Println("Starting Claude agent...")
		fmt.Println("Using prompt:\n=====\n", prompt)
		fmt.Println("=====")

		cmd := exec.Command("claude", "-p", prompt)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		
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
