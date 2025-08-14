
package cmd

import (
  "fmt"
)

// RunPipeline spawns multiple Claude agents that will perform different steps of a task
// Right now the only available pipeline is: Analyze -> Develop -> Review
func ClaudePipeline(branch string) {
  // TODO: Wait for user input before continuing after each step (i.e human-in-the-loop)
  // TODO: add a flag to skip user input so it runs full async

  fmt.Println("Launching Astropath's Pipeline of agents...")

  // Step 1: Analyze
  fmt.Println("Pipeline - Step #1. Analyze...")
  ClaudeAnalyze()

  // Step 2: Develop
  fmt.Println("Pipeline - Step #2. Develop...")
  ClaudeDevelop(branch)

  // Step 3: Review
  fmt.Println("Pipeline - Step #3. Review...")
  ClaudeReview(branch)
}
