package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fynardo/astropath/internal/claude"
)

func ClaudeExplore() {
	fmt.Println("Launching Claude explorer agent...")
	prompt := claude.GetPrompt(claude.ExplorerPrompt)
	done := claude.RunAgent(prompt)
	
	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Claude explorer agent launched. Use Ctrl+C to stop.")
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Claude agent exited with error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Claude agent completed successfully.")
	}
}
