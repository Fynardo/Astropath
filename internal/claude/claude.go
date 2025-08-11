package claude

import (
	"fmt"
	"os"
	"os/exec"
)

// Package claude provides a way for interacting with Claude Code.

const basePrompt = `You are a helpful AI assistant. Please help the user with their software engineering tasks.
Focus on providing clear, actionable solutions and follow best practices.

Always update the file ./ASTROPATH.md with your feedback, but don't overwrite it, append a section at the end.
It is Markdown, so start with a '#' to create a title, and then add your text inside that section.`

const defaultPrompt = basePrompt

const explorerPrompt = basePrompt + "\n" + `Your next task is to explore a project.
	Please take a look at the current dir (and subdirs) and identify:
	1. What the project is about
	2. Which technology stack is used
	3. What are the main components of the project
	Keep it short, don't think too much, just do a basic exploration.

	Don't forget to add your findings to the ./ASTROPATH.md file.
`

const reviewerPrompt = basePrompt + "\n" + `For your next task you are going to be a code reviewer AI assistant.
	1. Please checkout to the branch: {{ .BranchName }}.
	2. Get a diff of the latest commit: 'git show'
	3. Review the code and provide feedback.

	Don't forget to add your findings to the ./ASTROPATH.md file.
`

type PromptParams struct {
}

type ReviewerParams struct {
	BranchName string // Name of the pull request to review
}


// PromptType represents different types of prompts available
type PromptType string

const (
	DefaultPrompt PromptType = "default"
	ExplorerPrompt PromptType = "explorer"
	ReviewerPrompt PromptType = "reviewer"
)

// getPrompt returns the appropriate prompt based on the prompt type
func GetPrompt(promptType PromptType) string {
	switch promptType {
	case ExplorerPrompt:
		return explorerPrompt
	case ReviewerPrompt:
		return reviewerPrompt
	default:
		return defaultPrompt
	}
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
