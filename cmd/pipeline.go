
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// pipelineCmd represents the pipeline command
var pipelineCmd = &cobra.Command{
	Use:   "pipeline [branch]",
	Short: "Launch a Pipeline of Claude agents that will solve a task",
	Long: `Launch a Pipeline of Claude agents that will solve a task by analyzing -> developing -> reviewing it.

The pipeline executes three steps in sequence:
1. Analyze - Analyzes the issue and proposes a solution
2. Develop - Implements the proposed solution
3. Review - Reviews the implementation

If no branch is provided, the develop and review steps will use appropriate defaults.

Examples:
  astropath pipeline
  astropath pipeline my-feature-branch`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var branch string
		if len(args) > 0 {
			branch = args[0]
		}
		return claudePipeline(cmd, branch)
	},
}

func claudePipeline(cmd *cobra.Command, branch string) error {
	// TODO: Wait for user input before continuing after each step (i.e human-in-the-loop)
	// TODO: add a flag to skip user input so it runs full async

	fmt.Println("Launching Astropath's Pipeline of agents...")

	// Step 1: Analyze
	fmt.Println("Pipeline - Step #1. Analyze...")
	if err := claudeAnalyze(cmd); err != nil {
		return fmt.Errorf("pipeline step 1 (analyze) failed: %v", err)
	}

	// Step 2: Develop
	fmt.Println("Pipeline - Step #2. Develop...")
	if err := claudeDevelop(cmd, branch); err != nil {
		return fmt.Errorf("pipeline step 2 (develop) failed: %v", err)
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

// Legacy function for backward compatibility
func ClaudePipeline(branch string) {
	if err := claudePipeline(nil, branch); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
