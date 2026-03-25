package main

import (
	"fmt"
	"os"

	"github.com/vorpalengineering/cli/internal/commands"
	"github.com/vorpalengineering/cli/internal/commands/config"
	"github.com/vorpalengineering/cli/internal/commands/knowledge"
)

func main() {
	if len(os.Args) < 2 {
		commands.PrintHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "config":
		if len(os.Args) <= 2 {
			config.View()
			return
		}
		switch os.Args[2] {
		case "set":
			config.Set(os.Args[3:])
		default:
			fmt.Fprintf(os.Stderr, "Unknown config command: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "knowledge":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: vorpal knowledge <search|list>")
			os.Exit(1)
		}
		switch os.Args[2] {
		case "search":
			knowledge.Search(os.Args[3:])
		case "list":
			knowledge.List(os.Args[3:])
		default:
			fmt.Fprintf(os.Stderr, "Unknown knowledge command: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "version":
		commands.PrintVersion()

	case "help", "--help", "-h":
		commands.PrintHelp()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		commands.PrintHelp()
		os.Exit(1)
	}
}
