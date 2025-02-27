package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/DNelson35/JumpDir/search"
)

func printUsage() {
	fmt.Println("Usage: go run main.go <target_directory> [<starting_point>]")
	fmt.Println()
	fmt.Println("Arguments:")
	fmt.Println("  <target_directory>    Target directory (required)")
	fmt.Println("  [<starting_point>]    Starting point (optional)")
	fmt.Println()
	fmt.Println("Flags:")
	flag.PrintDefaults()
}

func main() {
	var help bool
	flag.BoolVar(&help, "help", false, "Show help information")
	flag.BoolVar(&help, "h", false, "Show help information")
	flag.Parse()

	if help {
		printUsage()
		os.Exit(0)
	}

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Error: <target_directory> is required.")
		printUsage()
		os.Exit(1)
	}

	// Note: this is trying to make the starting point optional.
	// (since it said in the usage that it would be optional)
	if len(args) < 2 {
		// shoutout to: https://stackoverflow.com/a/18537419/27913035
		ex, err := os.Executable()
		if err != nil {
			fmt.Println("an error occurred")
			fmt.Println(err.Error())
			os.Exit(1)
		}

		args = append(args, filepath.Dir(ex))
	}

	name := args[0]
	currDir := args[1]

	result := search.JumpDirectory(name, currDir)
	fmt.Println(result)
}
