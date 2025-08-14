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

type DeveloperParams struct {
	BranchName string
}

// developCmd represents the develop command
var developCmd = &cobra.Command{
	Use:   "develop [branch]",
	Short: "Launch a Claude agent that will write code as a developer",
	Long: `Launch a Claude agent that will write code as a developer.

If no branch name is provided, the agent will find a suitable name for it.

Examples:
  astropath develop
  astropath develop my-feature-branch`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var branch string
		if len(args) > 0 {
			branch = args[0]
		}
		return claudeDevelop(cmd, branch)
	},
}

func claudeDevelop(cmd *cobra.Command, branch string) error {
	fmt.Println("Launching Astropath's Claude Developer agent...")

	prompt := config.GetPrompt(config.DeveloperPromptType)
	promptParams := DeveloperParams{BranchName: branch}

	templ, err := template.New("prompt").Parse(prompt)
	if err != nil {
		return fmt.Errorf("parsing prompt template: %v", err)
	}

	var buff = bytes.Buffer{}
	err = templ.Execute(&buff, promptParams)
	if err != nil {
		return fmt.Errorf("executing prompt template: %v", err)
	}

	// Check if streaming flag is set, default to true for develop
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
		fmt.Println("Astropath's Claude Developer agent launched with streaming. Use Ctrl+C to stop.")
	} else {
		fmt.Println("Astropath's Claude Developer agent launched. Use Ctrl+C to stop.")
	}
	
	// Wait for the agent to complete
	select {
	case err := <-done:
		if err != nil {
			return fmt.Errorf("Claude Developer agent exited with error: %v", err)
		}
		fmt.Println("Claude agent completed successfully.")
		return nil
	}
}
