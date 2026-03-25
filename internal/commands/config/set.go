package config

import (
	"flag"
	"fmt"
	"os"
	"strings"

	appconfig "github.com/vorpalengineering/cli/internal/config"
)

func Set(args []string) {
	fs := flag.NewFlagSet("config set", flag.ExitOnError)
	apiKey := fs.String("api-key", "", "API key")
	apiURL := fs.String("api-url", "", "API base URL")
	fs.Parse(args)

	if *apiKey == "" && *apiURL == "" {
		fmt.Fprintln(os.Stderr, "Usage: vorpal config set --api-key <key> [--api-url <url>]")
		os.Exit(1)
	}

	cfg, _ := appconfig.Load()

	if *apiKey != "" {
		if !strings.HasPrefix(*apiKey, "vk_live_") {
			fmt.Fprintln(os.Stderr, "Error: API key must start with 'vk_live_'")
			os.Exit(1)
		}
		cfg.APIKey = *apiKey
	}

	if *apiURL != "" {
		cfg.APIURL = *apiURL
	}

	if err := appconfig.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Vorpal CLI Config Updated")
}
