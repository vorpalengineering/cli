package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

type nodeTypeEntry struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type edgeTypeEntry struct {
	Name        string `json:"name"`
	IsSymmetric bool   `json:"is_symmetric"`
	Description string `json:"description"`
}

func Types(args []string) {
	fs := flag.NewFlagSet("knowledge types", flag.ExitOnError)
	nodesFlag := fs.Bool("nodes", false, "only show node types")
	edgesFlag := fs.Bool("edges", false, "only show edge types")
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	// If neither flag set, show both
	showNodes := *nodesFlag || (!*nodesFlag && !*edgesFlag)
	showEdges := *edgesFlag || (!*nodesFlag && !*edgesFlag)

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var nodeTypes []nodeTypeEntry
	var edgeTypes []edgeTypeEntry

	if showNodes {
		body, err := c.Get("/knowledge/node-types")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		var resp struct {
			Types []nodeTypeEntry `json:"types"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing node types: %v\n", err)
			os.Exit(1)
		}
		nodeTypes = resp.Types
	}

	if showEdges {
		body, err := c.Get("/knowledge/edge-types")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		var resp struct {
			Types []edgeTypeEntry `json:"types"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing edge types: %v\n", err)
			os.Exit(1)
		}
		edgeTypes = resp.Types
	}

	if *jsonOut {
		out := map[string]interface{}{}
		if showNodes {
			out["node_types"] = nodeTypes
		}
		if showEdges {
			out["edge_types"] = edgeTypes
		}
		data, _ := json.MarshalIndent(out, "", "  ")
		fmt.Println(string(data))
		return
	}

	if showNodes {
		if len(nodeTypes) == 0 {
			fmt.Println("No node types found.")
		} else {
			fmt.Println("Node Types:")
			for _, t := range nodeTypes {
				if t.Description != "" {
					fmt.Printf("  %-15s %s\n", t.Name, t.Description)
				} else {
					fmt.Printf("  %s\n", t.Name)
				}
			}
		}
	}

	if showEdges {
		if showNodes {
			fmt.Println()
		}
		if len(edgeTypes) == 0 {
			fmt.Println("No edge types found.")
		} else {
			fmt.Println("Edge Types:")
			for _, t := range edgeTypes {
				sym := ""
				if t.IsSymmetric {
					sym = " (symmetric)"
				}
				if t.Description != "" {
					fmt.Printf("  %-15s %s%s\n", t.Name, t.Description, sym)
				} else {
					fmt.Printf("  %s%s\n", t.Name, sym)
				}
			}
		}
	}
}
