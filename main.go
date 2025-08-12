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
	case "analyze":
		handleClaudeAnalyze()
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
	fmt.Println("    analyze  Launch a Claude agent that will analyze an issue and propse a solution.")
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

func handleClaudeAnalyze() {
	cmd.ClaudeAnalyze()
}

func handleClaudeExplore() {
	cmd.ClaudeExplore()
}

func handleClaudeReview(args []string) {
	if len(args) == 0 {
		cmd.ClaudeReview("main") // Default to main branch if no branch is specified
	}

	cmd.ClaudeReview(args[0])
}
