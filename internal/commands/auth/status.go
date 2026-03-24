package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/vorpalengineering/cli/internal/client"
	"github.com/vorpalengineering/cli/internal/config"
)

func Status() {
	cfg, _ := config.Load()

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
