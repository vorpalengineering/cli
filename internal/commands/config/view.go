package config

import (
	"fmt"
	"strings"

	appconfig "github.com/vorpalengineering/cli/internal/config"
)

func View() {
	cfg, _ := appconfig.Load()

	fmt.Printf("  Config:   %s\n", appconfig.Path())
	fmt.Printf("  API URL:  %s\n", cfg.APIURL)

	if cfg.APIKey == "" {
		fmt.Println("  API Key:  not configured")
		fmt.Println("\nRun: vorpal config set-key <your-key>")
	} else {
		prefix := cfg.APIKey[:16]
		masked := prefix + strings.Repeat("•", len(cfg.APIKey)-16)
		fmt.Printf("  API Key:  %s\n", masked)
	}
}
