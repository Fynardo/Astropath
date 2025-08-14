package cmd

import (
	"bytes"
	"fmt"
	"text/template"
	"time"

	"github.com/fynardo/astropath/internal/claude"
	"github.com/fynardo/astropath/config"
	"github.com/spf13/cobra"
)

// reviewCmd represents the review command
var reviewCmd = &cobra.Command{
	Use:   "review [branch]",
	Short: "Launch a Claude agent that will review a branch",
	Long: `Launch a Claude agent that will review a branch.

Defaults to 'main' if no branch is specified.

Examples:
  astropath review
  astropath review main
  astropath review feature-branch`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branch := "main" // Default to main branch
		if len(args) > 0 {
			branch = args[0]
		}
		return claudeReview(cmd, branch)
	},
}

func claudeReview(cmd *cobra.Command, branch string) error {
	fmt.Println("Launching Astropath's Claude reviewer agent...")

	prompt := config.GetPrompt(config.ReviewerPromptType)
	promptParams := claude.ReviewerParams{BranchName: branch}

	templ, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return fmt.Errorf("parsing prompt template: %v", err)
	}

	var buff = bytes.Buffer{}
	err = templ.Execute(&buff, promptParams)
	if err != nil {
		return fmt.Errorf("executing prompt template: %v", err)
	}

	// Check if streaming flag is set, default to true for review
	useStreaming := streaming || true
	
	var done <-chan error
	if useStreaming {
		done = claude.RunAgentWithStreaming(buff.String())
	} else {
		done = claude.RunAgent(buff.String())
	}
	
	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	if useStreaming {
		fmt.Println("Astropath's Claude Reviewer agent launched with streaming. Use Ctrl+C to stop.")
	} else {
		fmt.Println("Astropath's Claude Reviewer agent launched. Use Ctrl+C to stop.")
	}
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Claude Reviewer agent exited with error: %v", err)
		}
		fmt.Println("Claude agent completed successfully.")
		return nil
	}
}
