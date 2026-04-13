# Vorpal CLI

Command-line tool for the Vorpal Engineering platform.

## Install

```bash
brew install vorpalengineering/tap/vorpal
```

Or via the install script:

```bash
curl -sSfL https://raw.githubusercontent.com/vorpalengineering/cli/main/install.sh | sh
```

Or via Go:

```bash
go install github.com/vorpalengineering/cli/cmd/vorpal@latest
```

Or build from source:

```bash
go build -o vorpal ./cmd/vorpal
```

## Setup

1. Create an account at [console.vorpalengineering.com](https://console.vorpalengineering.com)
2. Generate an API key on the Access page
3. Configure the CLI:

```bash
vorpal config set --api-key <your-key>
```

## Command Tree

```
vorpal
├── config                          View current configuration
│   └── set                         Set configuration values
│       ├── --api-key <key>         API key
│       └── --api-url <url>         API base URL
│
├── knowledge
│   ├── search <text>               Search the knowledge base (keyword by default)
│   │   ├── --mode <mode>           Search mode: keyword or semantic (default keyword)
│   │   ├── --type <name>           Filter by node type
│   │   ├── --limit N               Max results (default 5, max 20)
│   │   ├── --threshold N           Similarity threshold for semantic (default 0.5)
│   │   └── --json                  Output as JSON
│   │
│   ├── list                        List knowledge nodes
│   │   ├── --limit N               Nodes per page (default 10)
│   │   ├── --offset N              Skip N nodes
│   │   ├── --type <name>           Filter by node type
│   │   └── --json                  Output as JSON
│   │
│   ├── get <id>                    Show a knowledge node with its relations
│   │   └── --json                  Output as JSON
│   │
│   ├── traverse <id>               Walk the knowledge graph from a starting node
│   │   ├── --depth N               Traversal depth, 1-5 (default 2)
│   │   └── --json                  Output as JSON
│   │
│   └── types                       List available node and edge types
│       ├── --nodes                 Only show node types
│       ├── --edges                 Only show edge types
│       └── --json                  Output as JSON
│
├── version                         Show CLI version
└── help                            Show help
```

## Usage

```bash
# View config
vorpal config

# Set API key
vorpal config set --api-key ve_live_...

# Search knowledge base (keyword by default)
vorpal knowledge search "reentrancy vulnerability"
vorpal knowledge search --json --limit 3 "oracle manipulation"

# Semantic search (embedding-based)
vorpal knowledge search --mode semantic "reentrancy vulnerability"
vorpal knowledge search --mode semantic --threshold 0.3 "access control"

# List knowledge nodes
vorpal knowledge list
vorpal knowledge list --type finding --limit 20

# Get a specific node
vorpal knowledge get <node-id>
vorpal knowledge get --json <node-id>

# Traverse the knowledge graph from a node
vorpal knowledge traverse <node-id>
vorpal knowledge traverse --depth 3 <node-id>
vorpal knowledge traverse --json --depth 1 <node-id>
```

## Configuration

Config is stored at `~/.vorpal/config.json`:

```json
{
    "api_key": "ve_live_...",
    "api_url": "https://api.vorpalengineering.com"
}
```
