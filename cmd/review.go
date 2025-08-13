package cmd

import (
  "bytes"
  "fmt"
  "os"
	"text/template"
  "time"

  "github.com/fynardo/astropath/internal/claude"
  "github.com/fynardo/astropath/config"
)

func ClaudeReview(branch string) {
  fmt.Println("Launching Astropath's Claude reviewer agent...")

	prompt := config.GetPrompt(config.ReviewerPromptType)
	promptParams := claude.ReviewerParams{BranchName: branch}

	templ, err := template.New("prompt").Parse(prompt)
	var buff = bytes.Buffer{}

	err = templ.Execute(&buff, promptParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing prompt template: %v\n", err)
		panic(err)
	}

	done := claude.RunAgent(buff.String())
	
	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Astropath's Claude Reviewer agent launched. Use Ctrl+C to stop.")
	
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
