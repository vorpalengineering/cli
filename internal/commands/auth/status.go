package auth

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

func Status(args []string) {
	fs := flag.NewFlagSet("auth status", flag.ExitOnError)
	jsonOut := fs.Bool("json", false, "output as JSON")
	fs.Parse(args)

	cfg, _ := config.Load()

	if *jsonOut {
		prefix := ""
		if len(cfg.APIKey) >= 16 {
			prefix = cfg.APIKey[:16]
		}
		connected := false
		if cfg.APIKey != "" {
			c, err := client.New(cfg)
			if err == nil {
				connected = c.Ping() == nil
			}
		}
		out, _ := json.MarshalIndent(map[string]interface{}{
			"api_key_prefix": prefix,
			"api_url":        cfg.APIURL,
			"connected":      connected,
		}, "", "  ")
		fmt.Println(string(out))
		return
	}

	fmt.Printf("  API URL:  %s\n", cfg.APIURL)

	if cfg.APIKey == "" {
		fmt.Println("  API Key:  not configured")
		fmt.Println("  Status:   not authenticated")
		fmt.Println("\nRun: vellma auth set-key <your-key>")
		return
	}

	prefix := cfg.APIKey[:16]
	masked := prefix + strings.Repeat("•", len(cfg.APIKey)-16)
	fmt.Printf("  API Key:  %s\n", masked)

	c, err := client.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Status:   %v\n", err)
		return
	}

	if err := c.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "  Status:   %v\n", err)
	} else {
		fmt.Println("  Status:   connected")
	}
}
