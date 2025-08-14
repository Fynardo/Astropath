package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fynardo/astropath/internal/claude"
	"github.com/fynardo/astropath/config"
	"github.com/spf13/cobra"
)

// exploreCmd represents the explore command
var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details",
	Long:  "Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details.",
	RunE: func(cmd *cobra.Command, args []string) error {
		return claudeExplore(cmd)
	},
}

func claudeExplore(cmd *cobra.Command) error {
	fmt.Println("Launching Claude explorer agent...")
	prompt := config.GetPrompt(config.ExplorerPromptType)
	
	// Check if streaming flag is set, default to false for explore
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
		fmt.Println("Claude explorer agent launched with streaming. Use Ctrl+C to stop.")
	} else {
		fmt.Println("Claude explorer agent launched. Use Ctrl+C to stop.")
	}
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Claude explorer agent exited with error: %v", err)
		}
		fmt.Println("Claude agent completed successfully.")
		return nil
	}
}

// Legacy function for backward compatibility
func ClaudeExplore() {
	if err := claudeExplore(nil); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
