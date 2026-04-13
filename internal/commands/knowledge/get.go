package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

type nodeDetail struct {
	ID        string                 `json:"id"`
	NodeType  string                 `json:"node_type"`
	Title     string                 `json:"title"`
	Details   string                 `json:"details"`
	Citations *string                `json:"citations"`
	Quality   int                    `json:"quality"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt string                 `json:"created_at"`
	Related   []relatedNode          `json:"related"`
}

type relatedNode struct {
	ID        string `json:"id"`
	NodeType  string `json:"node_type"`
	EdgeType  string `json:"edge_type"`
	Quality   int    `json:"quality"`
	Title     string `json:"title"`
	Direction string `json:"direction"`
}

func Get(args []string) {
	fs := flag.NewFlagSet("knowledge get", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: vorpal knowledge get <id> [--json]")
		os.Exit(1)
	}

	id := fs.Arg(0)

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	body, err := c.Get("/knowledge/" + id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		printJSON(body)
		return
	}

	var node nodeDetail
	if err := json.Unmarshal(body, &node); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("  ┌─ %s\n", node.Title)
	fmt.Printf("  │  Type: %s  Quality: %d/5\n", node.NodeType, node.Quality)
	fmt.Printf("  │  ID: %s\n", node.ID)
	fmt.Printf("  │\n")

	if node.Details != "" {
		for _, line := range wrapText(node.Details, 72) {
			fmt.Printf("  │  %s\n", line)
		}
	}

	if node.Citations != nil && *node.Citations != "" {
		fmt.Printf("  │\n")
		fmt.Printf("  │  Citations:\n")
		// Citations are stored as JSON array of {label, url}
		var refs []struct {
			Label string `json:"label"`
			URL   string `json:"url"`
		}
		if err := json.Unmarshal([]byte(*node.Citations), &refs); err == nil {
			for _, ref := range refs {
				fmt.Printf("  │    %s  %s\n", ref.Label, ref.URL)
			}
		}
	}

	if len(node.Related) > 0 {
		outbound := filterRelated(node.Related, "outbound")
		symmetric := filterRelated(node.Related, "symmetric")
		inbound := filterRelated(node.Related, "inbound")

		fmt.Printf("  │\n")
		fmt.Printf("  │  Related (%d):\n", len(node.Related))
		printRelatedGroup(outbound, "→")
		printRelatedGroup(symmetric, "↔")
		printRelatedGroup(inbound, "←")
	}

	fmt.Printf("  └─\n")
}

func filterRelated(nodes []relatedNode, direction string) []relatedNode {
	var out []relatedNode
	for _, r := range nodes {
		if r.Direction == direction {
			out = append(out, r)
		}
	}
	return out
}

func printRelatedGroup(nodes []relatedNode, arrow string) {
	for _, r := range nodes {
		name := strings.ReplaceAll(strings.ReplaceAll(r.EdgeType, "_", " "), "-", " ")
		fmt.Printf("  │    %s %-14s %s (%s)\n", arrow, name, r.Title, r.NodeType)
	}
}
