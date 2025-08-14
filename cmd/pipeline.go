
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var noPause bool

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline [branch]",
	Short: "Launch a Pipeline of Claude agents that will solve a task",
	Long: `Launch a Pipeline of Claude agents that will solve a task by analyzing -> developing -> reviewing it.

The pipeline executes three steps in sequence:
1. Analyze - Analyzes the issue and proposes a solution
2. Develop - Implements the proposed solution
3. Review - Reviews the implementation

By default, the pipeline runs in interactive mode, pausing after each step for user confirmation.
Use the --no-pause flag to run all steps without interruption.

If no branch is provided, the develop and review steps will use appropriate defaults.

Examples:
  astropath pipeline                    (interactive mode)
  astropath pipeline --no-pause        (non-interactive mode)
  astropath pipeline my-branch         (interactive with branch)
  astropath pipeline my-branch --no-pause  (non-interactive with branch)`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var branch string
		if len(args) > 0 {
			branch = args[0]
		}
		return claudePipeline(cmd, branch)
	},
}

func init() {
	pipelineCmd.Flags().BoolVar(&noPause, "no-pause", false, "Skip user confirmation prompts between pipeline steps")
}

// waitForUserInput prompts the user to continue after completing a step
// Returns true if user wants to continue, false if they want to abort
func waitForUserInput(stepName string) bool {
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Printf("Step %s finished. Continue? (Y/n): ", stepName)
		
		input, err := reader.ReadString('\n')
		if err != nil {
			// Handle EOF (Ctrl+D) or other input errors
			fmt.Println("\nInput error or EOF detected. Aborting pipeline.")
			return false
		}
		
		// Trim whitespace and convert to lowercase
		input = strings.TrimSpace(strings.ToLower(input))
		
		// Handle empty input (Enter key) as "yes"
		if input == "" || input == "y" || input == "yes" {
			return true
		}
		
		if input == "n" || input == "no" {
			return false
		}
		
		// Invalid input, ask again
		fmt.Println("Please enter 'Y', 'n', or press Enter to continue.")
	}
}

func claudePipeline(cmd *cobra.Command, branch string) error {
	fmt.Println("Launching Astropath's Pipeline of agents...")

	// Step 1: Analyze
	fmt.Println("Pipeline - Step #1. Analyze...")
	if err := claudeAnalyze(cmd); err != nil {
		return fmt.Errorf("pipeline step 1 (analyze) failed: %v", err)
	}

	// Pause after Analyze step (unless --no-pause flag is set)
	if !noPause {
		if !waitForUserInput("Analyze") {
			fmt.Println("Pipeline aborted by user.")
			return nil
		}
	}

	// Step 2: Develop
	fmt.Println("Pipeline - Step #2. Develop...")
	if err := claudeDevelop(cmd, branch); err != nil {
		return fmt.Errorf("pipeline step 2 (develop) failed: %v", err)
	}

	// Pause after Develop step (unless --no-pause flag is set)
	if !noPause {
		if !waitForUserInput("Develop") {
			fmt.Println("Pipeline aborted by user.")
			return nil
		}
	}

	// Step 3: Review
	fmt.Println("Pipeline - Step #3. Review...")
	if branch == "" {
		branch = "main" // Default for review step
	}
	if err := claudeReview(cmd, branch); err != nil {
		return fmt.Errorf("pipeline step 3 (review) failed: %v", err)
	}

	fmt.Println("Pipeline completed successfully!")
	return nil
}
