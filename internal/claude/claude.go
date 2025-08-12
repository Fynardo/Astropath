package claude

import (
	"fmt"
	"os"
	"os/exec"
)

// Package claude provides a way for interacting with Claude Code.

const basePrompt = `You are a helpful AI assistant. Please help the user with their software engineering tasks.
Focus on providing clear, actionable solutions and follow best practices.

Always update the file ./ASTROPATH.md with your feedback, but don't overwrite it from scratch.
What you are going to do is to edit a specific section of the file. It is Markdown, so identify sections
as blocks that start with a '#', and then add your text inside that section, the specific section title you need to
update will be provided to you in the task description in the following paragraphs.` + "\n"

const defaultPrompt = basePrompt

const explorerPrompt = basePrompt + "\n" + `Your next task is to explore a project.
	Please take a look at the current dir (and subdirs) and identify:
	1. What the project is about
	2. Which technology stack is used
	3. What are the main components of the project
	Keep it short, don't think too much, just do a basic exploration.

	Don't forget to edit the ./ASTROPATH.md file with your findings, use the section called 'Exploration Report'.
`

const reviewerPrompt = basePrompt + "\n" + `For your next task you are going to be a code reviewer AI assistant.
	You are going to review the udpates to the code in a branch, probably part of a pull request, so you will:
	1. Get a diff of the branch compared to main: 'git diff main {{ .BranchName }}'
	3. Review the code update and provide feedback.

	Don't forget to add your findings to the ./ASTROPATH.md file, your section is called 'Issue Explanation'.
`

const analystPrompt = basePrompt + "\n" + `For your next task you are going to be a software analyst AI assistant.
	You are going to review an Issue detailed in the ./ASTROPATH.md file, under the 'Issue Explanation' section of the file.
	Your task is to propose a solution that consists of:
	1. A list of bullet points explaining what you want to achieve
	2. A TO-DO list explaining how you would do it

	Always remember that you are an analyst, you don't write code, your task it to
	propose a high-level solution to the problem that a coder can implement.

	Don't forget to add your findings to the ./ASTROPATH.md file, your section is called 'Solution Proposal'.`


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
	AnalystPrompt PromptType = "analyst"
)

// getPrompt returns the appropriate prompt based on the prompt type
func GetPrompt(promptType PromptType) string {
	switch promptType {
	case AnalystPrompt:
		return analystPrompt
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
