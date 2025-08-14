package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fynardo/astropath/internal/claude"
	"github.com/fynardo/astropath/config"
	"github.com/spf13/cobra"
)

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Launch a Claude agent that will analyze an issue and propose a solution",
	Long:  "Launch a Claude agent that will analyze an issue and propose a solution.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return claudeAnalyze(cmd)
	},
}

func claudeAnalyze(cmd *cobra.Command) error {
	fmt.Println("Launching Astropath's Claude Analyst agent...")

	prompt := config.GetPrompt(config.AnalystPromptType)
	
	// Check if streaming flag is set, default to false for analyze
	useStreaming := streaming || false
	
	var done <-chan error
	if useStreaming {
		done = claude.RunAgentWithStreaming(prompt)
	} else {
		done = claude.RunAgent(prompt)
	}

	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	if useStreaming {
		fmt.Println("Astropath's Claude Analyst agent launched with streaming. Use Ctrl+C to stop.")
	} else {
		fmt.Println("Astropath's Claude Analyst agent launched. Use Ctrl+C to stop.")
	}
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Claude Analyst agent exited with error: %v", err)
		}
		fmt.Println("Claude agent completed successfully.")
		return nil
	}
}

// Legacy function for backward compatibility
func ClaudeAnalyze() {
	if err := claudeAnalyze(nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
