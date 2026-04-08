package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

type knowledgeEntry struct {
	ID       string `json:"id"`
	NodeType string `json:"node_type"`
	Title    string `json:"title"`
	Preview  string `json:"preview"`
	Quality  int    `json:"quality"`
}

func List(args []string) {
	fs := flag.NewFlagSet("knowledge list", flag.ExitOnError)
	limit := fs.Int("limit", 10, "entries per page")
	offset := fs.Int("offset", 0, "skip N entries")
	nodeType := fs.String("type", "", "filter by node type")
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	params := url.Values{}
	params.Set("limit", strconv.Itoa(*limit))
	params.Set("offset", strconv.Itoa(*offset))
	if *nodeType != "" {
		params.Set("types", *nodeType)
	}

	body, err := c.Get("/knowledge/nodes?" + params.Encode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		var raw json.RawMessage
		json.Unmarshal(body, &raw)
		out, _ := json.MarshalIndent(raw, "", "  ")
		fmt.Println(string(out))
		return
	}

	var resp struct {
		Results []knowledgeEntry `json:"results"`
		Total   int              `json:"total"`
		Limit   int              `json:"limit"`
		Offset  int              `json:"offset"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Results) == 0 {
		fmt.Println("No entries found.")
		return
	}

	for _, e := range resp.Results {
		fmt.Printf("  [%s] %s (%s)\n", e.ID[:8], e.Title, e.NodeType)
	}

	fmt.Printf("\nShowing %d-%d of %d entries\n", resp.Offset+1, resp.Offset+len(resp.Results), resp.Total)
}
