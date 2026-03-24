package main

import (
	"fmt"
	"os"

	"github.com/vorpalengineering/cli/internal/commands"
	"github.com/vorpalengineering/cli/internal/commands/auth"
	"github.com/vorpalengineering/cli/internal/commands/knowledge"
)

func main() {
	if len(os.Args) < 2 {
		printHelp()
		os.Exit(0)
	}

	switch os.Args[1] {
	case "auth":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: vellma auth <set-key|status>")
			os.Exit(1)
		}
		switch os.Args[2] {
		case "set-key":
			auth.SetKey(os.Args[3:])
		case "status":
			auth.Status()
		default:
			fmt.Fprintf(os.Stderr, "Unknown auth command: %s\n", os.Args[2])
			os.Exit(1)
		}

	case "knowledge":
		if len(os.Args) < 3 {
			fmt.Fprintln(os.Stderr, "Usage: vellma knowledge <search|list>")
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
		printHelp()

	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", os.Args[1])
		printHelp()
		os.Exit(1)
	}
}

func printHelp() {
	fmt.Println(`vellma - Vorpal Engineering CLI

Usage:
  vellma <command> [subcommand] [flags]

Commands:
  auth set-key <key>     Store your API key
  auth status            Show authentication status

  knowledge search <text>  Search the knowledge base
    --limit N              Max results (default 5, max 20)
  knowledge list           List knowledge entries
    --limit N              Entries per page (default 10)
    --offset N             Skip N entries
    --category <name>      Filter by category
    --severity <level>     Filter by severity

  version                Show CLI version
  help                   Show this help`)
}
