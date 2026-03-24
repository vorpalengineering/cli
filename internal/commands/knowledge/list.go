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
	ID       string  `json:"id"`
	Title    string  `json:"title"`
	Category string  `json:"category"`
	Severity *string `json:"severity"`
	Quality  int     `json:"quality"`
}

func List(args []string) {
	fs := flag.NewFlagSet("knowledge list", flag.ExitOnError)
	limit := fs.Int("limit", 10, "entries per page")
	offset := fs.Int("offset", 0, "skip N entries")
	category := fs.String("category", "", "filter by category")
	severity := fs.String("severity", "", "filter by severity")
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
	if *category != "" {
		params.Set("category", *category)
	}
	if *severity != "" {
		params.Set("severity", *severity)
	}

	body, err := c.Get("/knowledge?" + params.Encode())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var resp struct {
		Entries []knowledgeEntry `json:"entries"`
		Total   int              `json:"total"`
		Limit   int              `json:"limit"`
		Offset  int              `json:"offset"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Entries) == 0 {
		fmt.Println("No entries found.")
		return
	}

	for _, e := range resp.Entries {
		severity := ""
		if e.Severity != nil {
			severity = *e.Severity
		}
		fmt.Printf("  [%s] %s (%s) — %s\n", e.ID[:8], e.Title, e.Category, severity)
	}

	fmt.Printf("\nShowing %d-%d of %d entries\n", resp.Offset+1, resp.Offset+len(resp.Entries), resp.Total)
}
