package main

import (
	"fmt"
	"os"

	"github.com/fynardo/astropath/cmd"
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
	case "echo":
		handleEcho(os.Args[2:])
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
	fmt.Println("    echo     Echo the provided arguments")
	fmt.Println("    init     Initialize Astropath and Claude settings in the current directory.")
	fmt.Println("    analyze  Launch a Claude agent that will analyze an issue and propse a solution.")
	fmt.Println("    develop  Launch a Claude agent that will write code as a developer.")
	fmt.Println("    explore  Launch a Claude agent that will explore the current dir to find what the project is about and provide basic details.")
	fmt.Println("    review   Launch a Claude agent that will review a branch.")
	fmt.Println()
	fmt.Println("Use 'astropath <COMMAND> --help' for more information about a specific command.")
}

func handleEcho(args []string) {
	if len(args) == 0 {
		fmt.Println("Usage: astropath echo <message>")
		return
	}

	for i, arg := range args {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(arg)
	}
	fmt.Println()
}

func handleInit() {
	// Create ASTROPATH.md file
	astropathContent :=
`# Exploration Report

# Issue Explanation

# Solution Proposal

# Implemented Code

# Solution Review
`

	err := os.WriteFile("ASTROPATH.md", []byte(astropathContent), 0644)
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
		claudeBaseJson := `{
  "permissions": {
  "allow": [
    "Edit",
    "Bash(./astropath echo:*)",
    "Bash(git commit:*)",
    "Bash(git checkout:*)"
  ],
  "deny": []
  }
}`

		err := os.WriteFile(settingsPath, []byte(claudeBaseJson), 0644)
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
