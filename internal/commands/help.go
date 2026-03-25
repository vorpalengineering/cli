package commands

import "fmt"

func PrintHelp() {
	fmt.Println(`vorpal - Vorpal Engineering CLI

Usage:
  vorpal <command> [subcommand] [flags]

Commands:
  config                 View current configuration
  config set             Set configuration values
    --api-key <key>        API key
    --api-url <url>        API base URL

  knowledge search <text>  Search the knowledge base
    --limit N              Max results (default 5, max 20)
    --threshold N          Similarity threshold (default 0.5, lower = stricter)
    --json                 Output as JSON
  knowledge list           List knowledge entries
    --limit N              Entries per page (default 10)
    --offset N             Skip N entries
    --category <name>      Filter by category
    --severity <level>     Filter by severity
    --json                 Output as JSON

  version                Show CLI version
  help                   Show this help`)
}
