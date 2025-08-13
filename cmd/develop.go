package cmd

import (
  "fmt"
  "os"
  "time"

  "github.com/fynardo/astropath/internal/claude"
  "github.com/fynardo/astropath/config"
)

func ClaudeDevelop() {
  fmt.Println("Launching Astropath's Claude Developer agent...")

	prompt := config.GetPrompt(config.DeveloperPromptType)
	done := claude.RunAgent(prompt)

	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Astropath's Claude Developer agent launched. Use Ctrl+C to stop.")
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			fmt.Fprintf(os.Stderr, "Astropath's Claude Developer agent exited with error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Claude agent completed successfully.")
	}
}
