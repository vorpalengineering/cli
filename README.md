# Vorpal CLI

Command-line tool for the Vorpal Engineering platform.

## Install

```bash
brew install vorpalengineering/tap/vorpal
```

Or build from source:

```bash
go build -o vorpal .
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
│   ├── search <text>               Search the knowledge base
│   │   ├── --limit N               Max results (default 5, max 20)
│   │   ├── --threshold N           Similarity threshold (default 0.5)
│   │   └── --json                  Output as JSON
│   │
│   └── list                        List knowledge entries
│       ├── --limit N               Entries per page (default 10)
│       ├── --offset N              Skip N entries
│       ├── --category <name>       Filter by category
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

# Search knowledge base
vorpal knowledge search "reentrancy vulnerability"
vorpal knowledge search --json --limit 3 "oracle manipulation"
vorpal knowledge search --threshold 0.3 "access control"

# List knowledge entries
vorpal knowledge list
vorpal knowledge list --category security
```

## Configuration

Config is stored at `~/.vorpal/config.json`:

```json
{
    "api_key": "ve_live_...",
    "api_url": "https://api.vorpalengineering.com"
}
```
