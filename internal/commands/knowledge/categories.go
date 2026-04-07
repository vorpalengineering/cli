package knowledge

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

func Types(args []string) {
	fs := flag.NewFlagSet("knowledge types", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	cfg, _ := config.Load()
	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	body, err := c.Get("/knowledge/node-types")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if *jsonOut {
		printJSON(body)
		return
	}

	var resp struct {
		Types []struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"types"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing response: %v\n", err)
		os.Exit(1)
	}

	if len(resp.Types) == 0 {
		fmt.Println("No node types found.")
		return
	}

	fmt.Println("Available node types:")
	for _, t := range resp.Types {
		if t.Description != "" {
			fmt.Printf("  %-15s %s\n", t.Name, t.Description)
		} else {
			fmt.Printf("  %s\n", t.Name)
		}
	}
}
