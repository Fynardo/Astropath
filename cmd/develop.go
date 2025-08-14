package cmd

import (
  "bytes"
  "fmt"
  "os"
	"text/template"
  "time"

  "github.com/fynardo/astropath/internal/claude"
  "github.com/fynardo/astropath/config"
)

type DeveloperParams struct {
	BranchName string
}


func ClaudeDevelop(branch string) {
  fmt.Println("Launching Astropath's Claude Developer agent...")

	prompt := config.GetPrompt(config.DeveloperPromptType)
	promptParams := DeveloperParams{BranchName: branch}

	templ, err := template.New("prompt").Parse(prompt)
	var buff = bytes.Buffer{}

	err = templ.Execute(&buff, promptParams)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing prompt template: %v\n", err)
		panic(err)
	}

	done := claude.RunAgentWithStreaming(buff.String())

	// Give the goroutine a moment to start before returning
	time.Sleep(100 * time.Millisecond)
	fmt.Println("Astropath's Claude Developer agent launched with streaming. Use Ctrl+C to stop.")
	
	// Wait for the agent to complete
	select {
	case result := <-done:
		if result.Error != nil {
			fmt.Fprintf(os.Stderr, "Astropath's Claude Developer agent exited with error: %v\n", result.Error)
			os.Exit(1)
		}
		fmt.Println("Claude agent completed successfully.")
		
		// Optional: Access to buffer for future processing
		buffer := result.Buffer.GetBuffer()
		fmt.Printf("Captured %d bytes of streaming output.\n", len(buffer))
		
		jsonLines := result.Buffer.GetJSONLines()
		fmt.Printf("Parsed %d JSON lines from stream.\n", len(jsonLines))
	}
}
