package main

import (
	"fmt"
	"strings"
	"os"

	"github.com/fynardo/astropath/cmd"
	"github.com/fynardo/astropath/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "astropath",
	Short: "CLI to orchestrate the execution of local AI agents",
	Long:  "Astropath - CLI to orchestrate the execution of local AI agents",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Astropath and Claude settings in the current directory",
	Long:  "Initialize Astropath and Claude settings in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		handleInit()
	},
}

var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Launch a Claude agent that will analyze an issue and propose a solution",
	Long:  "Launch a Claude agent that will analyze an issue and propose a solution.",
	Run: func(cmd *cobra.Command, args []string) {
		handleClaudeAnalyze()
	},
}

var developCmd = &cobra.Command{
	Use:   "develop",
	Short: "Launch a Claude agent that will write code as a developer",
	Long:  "Launch a Claude agent that will write code as a developer.",
	Run: func(cmd *cobra.Command, args []string) {
		handleClaudeDevelop()
	},
}

var exploreCmd = &cobra.Command{
	Use:   "explore",
	Short: "Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details",
	Long:  "Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details.",
	Run: func(cmd *cobra.Command, args []string) {
		handleClaudeExplore()
	},
}

var reviewCmd = &cobra.Command{
	Use:   "review [branch]",
	Short: "Launch a Claude agent that will review a branch",
	Long:  "Launch a Claude agent that will review a branch. Defaults to 'main' if no branch is specified.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleClaudeReview(args)
	},
}

var rawCmd = &cobra.Command{
	Use:   "raw <prompt>",
	Short: "Launch a Claude agent with a custom prompt as argument (UNSAFE)",
	Long:  "Launch a Claude agent with a custom prompt as argument (UNSAFE).",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		handleClaudeRaw(args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(developCmd)
	rootCmd.AddCommand(exploreCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(rawCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func handleInit() {
	// Create ASTROPATH.md file
	err := os.WriteFile("ASTROPATH.md", []byte(config.AstropathBaseTemplate), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating ASTROPATH.md: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Created ASTROPATH.md file in the current directory.\n")

	// Create .claude directory if it does not exist
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		err := os.Mkdir(".claude", 0755)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating .claude directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created .claude directory in the current directory.\n")
	} else {
		fmt.Printf(".claude directory already exists, skipping.\n")
	}

	// Create .claude/settings.json file if it does not exist
	settingsPath := ".claude/settings.json"
	if _, err := os.Stat(settingsPath); os.IsNotExist(err) {
		err := os.WriteFile(settingsPath, []byte(config.ClaudeSettingsJson), 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating .claude/settings.json: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created .claude/settings.json with default Astropath permissions.\n")
	} else {
		fmt.Printf(".claude/settings.json already exists, skipping.\n")
	}

	fmt.Println("Initialization complete!")
}

func handleClaudeAnalyze() {
	cmd.ClaudeAnalyze()
}

func handleClaudeDevelop() {
	cmd.ClaudeDevelop()
}

func handleClaudeExplore() {
	cmd.ClaudeExplore()
}

func handleClaudeReview(args []string) {
	if len(args) == 0 {
		cmd.ClaudeReview("main") // Default to main branch if no branch is specified
	} else {
		cmd.ClaudeReview(args[0])
	}
}

func handleClaudeRaw(args []string) {
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "No prompt provided for raw Claude agent.\n")
		os.Exit(1)
	}
	prompt := strings.Join(args, " ") // The idea is that it works either the prompt is quoted or not
	cmd.ClaudeRaw(prompt)
}
