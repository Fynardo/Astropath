package main

import (
	"fmt"
	"strings"
	"os"

	"github.com/fynardo/astropath/cmd"
	"github.com/fynardo/astropath/config"
)


func main() {
	if len(os.Args) < 2 {
		showHelp()
		return
	}

	command := os.Args[1]

	switch command {
	case "help", "--help", "-h":
		showHelp()
	case "init":
		handleInit()
	case "analyze":
		handleClaudeAnalyze()
	case "develop":
		handleClaudeDevelop()
	case "explore":
		handleClaudeExplore()
	case "review":
		handleClaudeReview(os.Args[2:])
	case "raw":
		handleClaudeRaw(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		fmt.Fprintf(os.Stderr, "Run 'astropath help' for usage information.\n")
		os.Exit(1)
	}
}

func showHelp() {
	fmt.Println("Astropath - CLI to orchestrate the execution of local AI agents")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("    astropath <COMMAND>")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("    help     Show this help message")
	fmt.Println("    init     Initialize Astropath and Claude settings in the current directory.")
	fmt.Println("    analyze  Launch a Claude agent that will analyze an issue and propse a solution.")
	fmt.Println("    develop  Launch a Claude agent that will write code as a developer.")
	fmt.Println("    explore  Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details.")
	fmt.Println("    review   Launch a Claude agent that will review a branch.")
	fmt.Println("    raw      Launch a Claude agent with a custom prompt as argument (UNSAFE).")
	fmt.Println()
	fmt.Println("Use 'astropath <COMMAND> --help' for more information about a specific command.")
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
