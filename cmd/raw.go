package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fynardo/astropath/internal/claude"
	"github.com/spf13/cobra"
)

// rawCmd represents the raw command
var rawCmd = &cobra.Command{
	Use:   "raw <prompt>",
	Short: "Launch a Claude agent with a custom prompt as argument (UNSAFE)",
	Long: `Launch a Claude agent with a custom prompt as argument (UNSAFE).

This command allows you to provide a custom prompt directly to Claude.
The prompt can be provided as multiple arguments or as a quoted string.

Examples:
  astropath raw "Help me debug this code"
  astropath raw Help me debug this code`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompt := strings.Join(args, " ")
		return claudeRaw(cmd, prompt)
	},
}

func claudeRaw(cmd *cobra.Command, prompt string) error {
	if prompt == "" {
		return fmt.Errorf("no prompt provided for raw Claude agent")
	}

	fmt.Println("Launching Claude Raw agent...")
	
	// Check if streaming flag is set, default to false for raw
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
		fmt.Println("Claude Raw agent launched with streaming. Use Ctrl+C to stop.")
	} else {
		fmt.Println("Claude Raw agent launched. Use Ctrl+C to stop.")
	}

	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Claude Raw agent exited with error: %v", err)
		}
		fmt.Println("Claude Raw agent completed successfully.")
		return nil
	}
}

// Legacy function for backward compatibility
func ClaudeRaw(prompt string) {
	if err := claudeRaw(nil, prompt); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

