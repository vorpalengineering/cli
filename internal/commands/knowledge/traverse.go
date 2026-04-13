package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

type graphNode struct {
	ID       string `json:"id"`
	NodeType string `json:"node_type"`
	Title    string `json:"title"`
	Quality  int    `json:"quality"`
}

type graphEdge struct {
	SourceID string `json:"source_id"`
	TargetID string `json:"target_id"`
	EdgeType string `json:"edge_type"`
	Quality  int    `json:"quality"`
}

func Traverse(args []string) {
	fs := flag.NewFlagSet("knowledge traverse", flag.ExitOnError)
	depth := fs.Int("depth", 2, "traversal depth (1-5)")
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Usage: vorpal knowledge traverse <id> [--depth N] [--json]")
		os.Exit(1)
	}

	id := fs.Arg(0)

	if *depth < 1 || *depth > 5 {
		fmt.Fprintln(os.Stderr, "Error: depth must be between 1 and 5")
		os.Exit(1)
	}

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	body, err := c.Get("/knowledge/graph/" + id + "?depth=" + strconv.Itoa(*depth))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		printJSON(body)
		return
	}

	var resp struct {
		Nodes []graphNode `json:"nodes"`
		Edges []graphEdge `json:"edges"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Nodes) == 0 {
		fmt.Println("No nodes found.")
		return
	}

	// Build adjacency info
	nodeMap := make(map[string]graphNode)
	for _, n := range resp.Nodes {
		nodeMap[n.ID] = n
	}

	// Find the root node
	root, ok := nodeMap[id]
	if !ok && len(resp.Nodes) > 0 {
		root = resp.Nodes[0]
	}

	// Build connections per node
	type connection struct {
		NodeID   string
		EdgeType string
		Outbound bool
	}
	connections := make(map[string][]connection)
	for _, e := range resp.Edges {
		connections[e.SourceID] = append(connections[e.SourceID], connection{NodeID: e.TargetID, EdgeType: e.EdgeType, Outbound: true})
		connections[e.TargetID] = append(connections[e.TargetID], connection{NodeID: e.SourceID, EdgeType: e.EdgeType, Outbound: false})
	}

	// Print root
	fmt.Printf("  ┌─ %s\n", root.Title)
	fmt.Printf("  │  Type: %s  Quality: %d/5\n", root.NodeType, root.Quality)
	fmt.Printf("  │  ID: %s\n", root.ID)
	fmt.Printf("  │\n")

	// Print connected nodes
	if conns, ok := connections[root.ID]; ok && len(conns) > 0 {
		fmt.Printf("  │  Connected (%d):\n", len(conns))
		for _, conn := range conns {
			n := nodeMap[conn.NodeID]
			arrow := "→"
			if !conn.Outbound {
				arrow = "←"
			}
			fmt.Printf("  │    %s %-14s %s (%s)\n", arrow, conn.EdgeType, n.Title, n.NodeType)
		}
	}

	// Print other nodes not directly connected to root (depth > 1)
	printed := map[string]bool{root.ID: true}
	for _, conn := range connections[root.ID] {
		printed[conn.NodeID] = true
	}

	var others []graphNode
	for _, n := range resp.Nodes {
		if !printed[n.ID] {
			others = append(others, n)
		}
	}

	if len(others) > 0 {
		fmt.Printf("  │\n")
		fmt.Printf("  │  Extended (%d more nodes):\n", len(others))
		for _, n := range others {
			fmt.Printf("  │    [%s] %s (%s)\n", n.ID[:8], n.Title, n.NodeType)
		}
	}

	fmt.Printf("  │\n")
	fmt.Printf("  │  Summary: %d nodes, %d edges, depth %d\n", len(resp.Nodes), len(resp.Edges), *depth)
	fmt.Printf("  └─\n")
}
