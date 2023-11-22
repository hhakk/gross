package main

import (
	"fmt"
	"os"

	"github.com/hhakk/gross/feed"
	"github.com/hhakk/gross/tui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("usage: %s <url-file>\n", os.Args[0])
		os.Exit(1)
	}
	path := os.Args[1]
	urls, err := feed.GetURLs(path)
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
