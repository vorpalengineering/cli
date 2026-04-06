package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

func Categories(args []string) {
	fs := flag.NewFlagSet("knowledge categories", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	body, err := c.Get("/knowledge/categories")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		printJSON(body)
		return
	}

	var resp struct {
		Categories []string `json:"categories"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Categories) == 0 {
		fmt.Println("No categories found.")
		return
	}

	fmt.Println("Available categories:")
	for _, cat := range resp.Categories {
		fmt.Printf("  %s\n", cat)
	}
}
