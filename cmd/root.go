package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fynardo/astropath/config"
	"github.com/spf13/cobra"
)

var streaming bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "astropath",
	Short: "CLI to orchestrate the execution of local AI agents",
	Long:  "Astropath - CLI to orchestrate the execution of local AI agents",
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Add persistent flag for streaming
	rootCmd.PersistentFlags().BoolVar(&streaming, "streaming", true, "Enable streaming output (overrides command defaults)")

	// Add all commands
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(developCmd)
	rootCmd.AddCommand(exploreCmd)
	rootCmd.AddCommand(reviewCmd)
	rootCmd.AddCommand(rawCmd)
	rootCmd.AddCommand(pipelineCmd)
	rootCmd.AddCommand(refreshCmd)
}

// initCmd handles the initialization of Astropath
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Astropath and Claude settings in the current directory",
	Long:  "Initialize Astropath and Claude settings in the current directory.",
	Run: func(cmd *cobra.Command, args []string) {
		if err := handleInit(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func handleInit() error {
	// Create ASTROPATH.md file
	err := os.WriteFile("ASTROPATH.md", []byte(config.AstropathBaseTemplate), 0644)
	if err != nil {
		return fmt.Errorf("creating ASTROPATH.md: %v", err)
	}
	fmt.Printf("Created ASTROPATH.md file in the current directory.\n")

	// Create .claude directory if it does not exist
	if _, err := os.Stat(".claude"); os.IsNotExist(err) {
		err := os.Mkdir(".claude", 0755)
		if err != nil {
			return fmt.Errorf("creating .claude directory: %v", err)
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
			return fmt.Errorf("creating .claude/settings.json: %v", err)
		}
		fmt.Printf("Created .claude/settings.json with default Astropath permissions.\n")
	} else {
		fmt.Printf(".claude/settings.json already exists, skipping.\n")
	}

	fmt.Println("Initialization complete!")
	return nil
}

// refreshCmd handles clearing ASTROPATH.md while preserving the Exploration Report
var refreshCmd = &cobra.Command{
	Use:   "refresh",
	Short: "Clear ASTROPATH.md while preserving the Exploration Report section",
	Long: `Refresh clears the ASTROPATH.md file to provide a clean slate for new tasks,
while preserving the valuable 'Exploration Report' section that contains project context.

This command will:
- Preserve the 'Exploration Report' section and its content
- Clear all other sections (Issue Explanation, Solution Proposal, Implemented Code, Code Review)
- Reset those sections to their empty template state

Use --force to skip the confirmation prompt.`,
	Run: func(cmd *cobra.Command, args []string) {
		force, _ := cmd.Flags().GetBool("force")
		if err := handleRefresh(force); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	},
}

func init() {
	refreshCmd.Flags().BoolP("force", "f", false, "Skip confirmation prompt")
}

func handleRefresh(force bool) error {
	// Check if ASTROPATH.md exists
	if _, err := os.Stat("ASTROPATH.md"); os.IsNotExist(err) {
		return fmt.Errorf("ASTROPATH.md does not exist. Run 'astropath init' first")
	}

	// Ask for confirmation unless --force is used
	if !force {
		fmt.Print("This will clear all sections except 'Exploration Report'. Are you sure? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("reading input: %v", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Read current content
	content, err := os.ReadFile("ASTROPATH.md")
	if err != nil {
		return fmt.Errorf("reading ASTROPATH.md: %v", err)
	}

	// Extract the Exploration Report section
	explorationReport := extractExplorationReport(string(content))

	// Create new content with preserved Exploration Report and fresh template sections
	newContent := explorationReport + "\n# Issue Explanation\n\n# Solution Proposal\n\n# Implemented Code\n\n# Code Review\n"

	// Write the refreshed content
	err = os.WriteFile("ASTROPATH.md", []byte(newContent), 0644)
	if err != nil {
		return fmt.Errorf("writing ASTROPATH.md: %v", err)
	}

	fmt.Println("ASTROPATH.md refreshed successfully!")
	fmt.Println("- Preserved: Exploration Report section")
	fmt.Println("- Cleared: Issue Explanation, Solution Proposal, Implemented Code, Code Review")
	return nil
}

func extractExplorationReport(content string) string {
	lines := strings.Split(content, "\n")
	var explorationLines []string
	inExplorationSection := false
	
	for _, line := range lines {
		// Check if we're starting the Exploration Report section
		if strings.HasPrefix(line, "# Exploration Report") {
			inExplorationSection = true
			explorationLines = append(explorationLines, line)
			continue
		}
		
		// Check if we've hit another section (any line starting with #)
		if inExplorationSection && strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "# Exploration Report") {
			break
		}
		
		// If we're in the exploration section, add the line
		if inExplorationSection {
			explorationLines = append(explorationLines, line)
		}
	}
	
	// If no exploration report was found, create a default one
	if len(explorationLines) == 0 {
		return "# Exploration Report\n"
	}
	
	return strings.Join(explorationLines, "\n")
}
