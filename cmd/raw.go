package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fynardo/astropath/internal/claude"
)

func ClaudeRaw(prompt string) {
	fmt.Println("Launching Claude Raw agent...")
	done := claude.RunAgentWithStreaming(prompt)

	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Claude Raw agent launched. Use Ctrl+C to stop.")

	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Claude Raw agent exited with error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Claude Raw agent completed successfully.")
	}
}

