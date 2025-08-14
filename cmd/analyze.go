package cmd

import (
  "fmt"
  "os"
  "time"

  "github.com/fynardo/astropath/internal/claude"
  "github.com/fynardo/astropath/config"
)

func ClaudeAnalyze() {
  fmt.Println("Launching Astropath's Claude Analyst agent...")

	prompt := config.GetPrompt(config.AnalystPromptType)
	done := claude.RunAgentWithStreaming(prompt)

	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Astropath's Claude Analyst agent launched with streaming. Use Ctrl+C to stop.")
	
	// Wait for the agent to complete
	select {
	case result := <-done:
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "Astropath's Claude Analyst agent exited with error: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("Claude agent completed successfully.")
		
		// Optional: Access to buffer for future processing
		buffer := result.Buffer.GetBuffer()
		fmt.Printf("Captured %d bytes of streaming output.\n", len(buffer))
		
		jsonLines := result.Buffer.GetJSONLines()
		fmt.Printf("Parsed %d JSON lines from stream.\n", len(jsonLines))
	}
}
