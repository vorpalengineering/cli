package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

type searchResult struct {
	Title        string  `json:"title"`
	Category     string  `json:"category"`
	Content      string  `json:"content"`
	CodeExamples *string `json:"code_examples"`
	Mitigation   *string `json:"mitigation"`
	Quality      int     `json:"quality"`
}

func Search(args []string) {
	fs := flag.NewFlagSet("knowledge search", flag.ExitOnError)
	limit := fs.Int("limit", 5, "max results (1-20)")
	mode := fs.String("mode", "keyword", "search mode: keyword or semantic")
	threshold := fs.Float64("threshold", 0.5, "similarity threshold for semantic search (0-2, lower = stricter)")
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	text := strings.Join(fs.Args(), " ")
	if text == "" {
		fmt.Fprintln(os.Stderr, "Usage: vorpal knowledge search <text> [--mode keyword|semantic] [--limit N] [--threshold N] [--json]")
		os.Exit(1)
	}

	if *mode != "keyword" && *mode != "semantic" {
		fmt.Fprintln(os.Stderr, "Error: --mode must be 'keyword' or 'semantic'")
		os.Exit(1)
	}

	// Warn if --threshold is used without semantic mode
	if *mode != "semantic" && *threshold != 0.5 {
		fmt.Fprintln(os.Stderr, "Note: --threshold is ignored without --mode semantic")
	}

	if *mode == "semantic" && (*threshold < 0 || *threshold > 2) {
		fmt.Fprintln(os.Stderr, "Error: threshold must be between 0 and 2")
		os.Exit(1)
	}

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	var results []searchResult

	if *mode == "semantic" {
		// Semantic search via POST /knowledge/search
		body, err := c.Post("/knowledge/search", map[string]interface{}{
			"text":      text,
			"limit":     *limit,
			"threshold": *threshold,
		})
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if *jsonOut {
			printJSON(body)
			return
		}

		var resp struct {
			Results []searchResult `json:"results"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}
		results = resp.Results
	} else {
		// Keyword search via GET /knowledge?search=...
		params := url.Values{}
		params.Set("search", text)
		params.Set("limit", strconv.Itoa(*limit))

		body, err := c.Get("/knowledge?" + params.Encode())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if *jsonOut {
			printJSON(body)
			return
		}

		var resp struct {
			Chunks []searchResult `json:"chunks"`
		}
		if err := json.Unmarshal(body, &resp); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
			os.Exit(1)
		}
		results = resp.Chunks
	}

	if len(results) == 0 {
		fmt.Println("No results found.")
		return
	}

	fmt.Printf("%d result(s) found:\n\n", len(results))

	for i, r := range results {
		printResult(i+1, r)
	}
}

func printResult(num int, r searchResult) {
	fmt.Printf("  ┌─ [%d] %s\n", num, r.Title)
	fmt.Printf("  │  Category: %s  Quality: %d/5\n", r.Category, r.Quality)
	fmt.Printf("  │\n")

	fmt.Printf("  │  Description:\n")
	for _, line := range wrapText(r.Content, 72) {
		fmt.Printf("  │    %s\n", line)
	}

	if r.CodeExamples != nil && *r.CodeExamples != "" {
		fmt.Printf("  │\n")
		fmt.Printf("  │  Code Examples:\n")
		for _, line := range strings.Split(*r.CodeExamples, "\n") {
			fmt.Printf("  │    %s\n", line)
		}
	}

	if r.Mitigation != nil && *r.Mitigation != "" {
		fmt.Printf("  │\n")
		fmt.Printf("  │  Mitigation:\n")
		for _, line := range wrapText(*r.Mitigation, 72) {
			fmt.Printf("  │    %s\n", line)
		}
	}

	fmt.Printf("  └─\n\n")
}

func printJSON(body []byte) {
	var raw json.RawMessage
	json.Unmarshal(body, &raw)
	out, _ := json.MarshalIndent(raw, "", "  ")
	fmt.Println(string(out))
}

func wrapText(s string, width int) []string {
	var lines []string
	words := strings.Fields(s)
	if len(words) == 0 {
		return lines
	}
	current := words[0]
	for _, w := range words[1:] {
		if len(current)+1+len(w) > width {
			lines = append(lines, current)
			current = w
		} else {
			current += " " + w
		}
	}
	if current != "" {
		lines = append(lines, current)
	}
	return lines
}
