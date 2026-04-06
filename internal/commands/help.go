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

  knowledge search <text>  Search the knowledge base (keyword by default)
    --mode <mode>          Search mode: keyword or semantic (default keyword)
    --limit N              Max results (default 5, max 20)
    --threshold N          Similarity threshold for semantic (default 0.5)
    --json                 Output as JSON
  knowledge list           List knowledge entries
    --limit N              Entries per page (default 10)
    --offset N             Skip N entries
    --category <name>      Filter by category
    --json                 Output as JSON
  knowledge categories     List available categories

  version                Show CLI version
  help                   Show this help`)
}
