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

type searchResult struct {
	Title      string  `json:"title"`
	Category   string  `json:"category"`
	Severity   *string `json:"severity"`
	Content    string  `json:"content"`
	Mitigation *string `json:"mitigation"`
	Quality    int     `json:"quality"`
}

func Search(args []string) {
	fs := flag.NewFlagSet("knowledge search", flag.ExitOnError)
	limit := fs.Int("limit", 5, "max results (1-20)")
	fs.Parse(args)

	text := strings.Join(fs.Args(), " ")
	if text == "" {
		fmt.Fprintln(os.Stderr, "Usage: vellma knowledge search <text> [--limit N]")
		os.Exit(1)
	}

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	body, err := c.Post("/knowledge/search", map[string]interface{}{
		"text":  text,
		"limit": *limit,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var resp struct {
		Results []searchResult `json:"results"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Results) == 0 {
		fmt.Println("No results found.")
		return
	}

	for i, r := range resp.Results {
		severity := "unknown"
		if r.Severity != nil {
			severity = *r.Severity
		}
		fmt.Printf("[%d] %s (%s — %s)\n", i+1, r.Title, r.Category, severity)
		fmt.Printf("    %s\n", truncate(r.Content, 200))
		if r.Mitigation != nil && *r.Mitigation != "" {
			fmt.Printf("    Mitigation: %s\n", truncate(*r.Mitigation, 150))
		}
		fmt.Println()
	}
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max] + "..."
}
