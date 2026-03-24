package commands

import "fmt"

const Version = "0.1.0"

func PrintVersion() {
	fmt.Printf("vellma CLI v%s\n", Version)
}
