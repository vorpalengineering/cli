package commands

import "fmt"

func PrintHelp() {
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
