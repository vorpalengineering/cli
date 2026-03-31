package commands

import "fmt"

const Version = "0.1.2"

func PrintVersion() {
	fmt.Printf("Vorpal Engineering CLI v%s\n", Version)
}
