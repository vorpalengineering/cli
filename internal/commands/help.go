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
    --type <name>          Filter by node type
    --limit N              Max results (default 5, max 20)
    --threshold N          Similarity threshold for semantic (default 0.5)
    --json                 Output as JSON
  knowledge list           List knowledge nodes
    --limit N              Nodes per page (default 10)
    --offset N             Skip N nodes
    --type <name>          Filter by node type
    --json                 Output as JSON
  knowledge get <id>       Show a knowledge node with its relations
    --json                 Output as JSON
  knowledge traverse <id>  Walk the knowledge graph from a starting node
    --depth N              Traversal depth, 1-5 (default 2)
    --json                 Output as JSON
  knowledge types          List available node and edge types
    --nodes                Only show node types
    --edges                Only show edge types
    --json                 Output as JSON

  version                Show CLI version
  help                   Show this help`)
}
