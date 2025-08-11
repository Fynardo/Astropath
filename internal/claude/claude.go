package claude

import (
	"fmt"
	"os"
	"os/exec"
)

// Package claude provides a way for interacting with Claude Code.

const defaultPrompt = `You are a helpful AI assistant. Please help the user with their software engineering tasks.
Focus on providing clear, actionable solutions and follow best practices.`

const explorerPrompt = `You are a coder assistant. 
	Please take a look at the current dir (and subdirs) and identify:
	1. What the project is about
	2. Which technology stack is used
	3. What are the main components of the project
	Keep it short, don't think too much, just do a basic exploration.
`

type PromptParams struct {
}


// PromptType represents different types of prompts available
type PromptType string

const (
	DefaultPrompt PromptType = "default"
	ExplorerPrompt PromptType = "explore"
)

// getPrompt returns the appropriate prompt based on the prompt type
func getPrompt(promptType PromptType) string {
	switch promptType {
	case ExplorerPrompt:
		return explorerPrompt
	default:
		return defaultPrompt
	}
}

// RunAgent spawns a Claude agent using the 'claude -p' command with the specified prompt type.
// It returns a channel that will receive a signal when the agent finishes.
func RunAgent(promptType PromptType, promptParams PromptParams ) <-chan error {
	done := make(chan error, 1)
	
	go func() {
		fmt.Println("Starting Claude agent...")
		
		prompt := getPrompt(promptType)
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
