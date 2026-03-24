package auth

import (
	"fmt"
	"os"
	"strings"

	"github.com/vorpalengineering/cli/internal/config"
)

func SetKey(args []string) {
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Usage: vellma auth set-key <api-key>")
		os.Exit(1)
	}

	key := args[0]
	if !strings.HasPrefix(key, "vk_live_") {
		fmt.Fprintln(os.Stderr, "Error: API key must start with 'vk_live_'")
		os.Exit(1)
	}

	cfg, _ := config.Load()
	cfg.APIKey = key
	if err := config.Save(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		os.Exit(1)
	}

	prefix := key[:16]
	fmt.Printf("API key saved: %s%s\n", prefix, strings.Repeat("•", len(key)-16))
}
