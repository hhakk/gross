package main

import (
	"fmt"
	"os"

	"github.com/hhakk/gross/feed"
	"github.com/hhakk/gross/tui"
)

func main() {
	path := os.Args[1]
	u, err := feed.GetURLs(path)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
	err = tui.Run(urls)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		os.Exit(1)
	}
}
