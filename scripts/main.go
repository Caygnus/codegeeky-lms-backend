package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/omkar273/codegeeky/scripts/internal"
)

// Command represents a script that can be run
type Command struct {
	Name        string
	Description string
	Run         func() error
}

var commands = []Command{
	{
		Name:        "generate-ent",
		Description: "Generate the Ent schema",
		Run:         internal.GenerateEnt,
	},
}

func main() {
	// Define command line flags
	var (
		listCommands bool
		cmdName      string
	)

	flag.BoolVar(&listCommands, "list", false, "List all available commands")
	flag.StringVar(&cmdName, "cmd", "", "Command to run")

	flag.Usage = func() {
		fmt.Printf("Usage of %s:\n\n", "scripts")
		fmt.Println("This tool helps you run various development and maintenance scripts.")
		fmt.Println("\nOptions:")
		flag.PrintDefaults()
		fmt.Println("\nExamples:")
		fmt.Println("  go run scripts/main.go -list                    # List all available commands")
		fmt.Println("  go run scripts/main.go -cmd generate-ent        # Generate Ent schema")
		fmt.Println("  go run scripts/main.go -cmd seed-ranks          # Seed rank data")
		fmt.Println("  go run scripts/main.go -cmd clear-ranks         # Clear all ranks")
		fmt.Println("  go run scripts/main.go -cmd validate-ranks      # Validate seeded ranks")
		fmt.Println("\nNote: Ensure your database configuration is properly set up before running database commands.")
	}

	flag.Parse()

	if listCommands {
		fmt.Println("Available commands:")
		fmt.Println()
		for _, cmd := range commands {
			fmt.Printf("  %-20s %s\n", cmd.Name, cmd.Description)
		}
		fmt.Println()
		fmt.Println("Use: go run scripts/main.go -cmd <command-name>")
		return
	}

	if cmdName == "" {
		fmt.Println("❌ Please specify a command to run using -cmd flag.")
		fmt.Println("💡 Use -list to see available commands or -help for usage information.")
		log.Fatal("No command specified")
	}

	// Find and run the command
	for _, cmd := range commands {
		if cmd.Name == cmdName {
			log.Printf("🚀 Running command: %s", cmdName)
			if err := cmd.Run(); err != nil {
				log.Fatalf("❌ Error running command %s: %v", cmdName, err)
			}
			log.Printf("✅ Command %s completed successfully!", cmdName)
			return
		}
	}

	log.Fatalf("❌ Unknown command: %s. Use -list to see available commands.", cmdName)
}
