package cmd

import (
	"fmt"
	"os"

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
