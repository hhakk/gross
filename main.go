package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hhakk/gross/feed"
	"github.com/hhakk/gross/tui"
	"github.com/spf13/viper"
)

func main() {
	configdir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("error reading config directory")
		os.Exit(1)
	}
	grossconfig := filepath.Join(configdir, "gross")
	err = os.MkdirAll(grossconfig, os.ModePerm)
	if err != nil {
		fmt.Println("error reading config directory")
		os.Exit(1)
	}
	cachedir, err := os.UserCacheDir()
	if err != nil {
		fmt.Println("error reading cache directory")
		os.Exit(1)
	}
	grosscache := filepath.Join(cachedir, "gross")
	err = os.MkdirAll(grosscache, os.ModePerm)
	if err != nil {
		fmt.Println("error reading cache directory")
		os.Exit(1)
	}
	browser, ok := os.LookupEnv("$BROWSER")
	fmt.Println(browser)
	if !ok {
		browser = "firefox"
	}
	// cachedir, err := os.UserCacheDir()
	viper.SetDefault("browsercmd", browser)
	viper.SetDefault("cachedir", grosscache)
	viper.SetDefault("urlfile", filepath.Join(grossconfig, "urls"))
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath(grossconfig)

	err = viper.ReadInConfig()
	if err != nil {
		fmt.Println("no config file found, using defaults.")
	}

	path := viper.GetString("urlfile")
	urls, err := feed.GetURLs(path)
	if err != nil || len(urls) < 1 {
		fmt.Printf("error: do you have a file with URLs in '%s'?\n", path)
		os.Exit(1)
	}
	err = tui.Run(urls)
	if err != nil {
		fmt.Printf("error: '%s'\n", err)
		os.Exit(1)
	}
}
