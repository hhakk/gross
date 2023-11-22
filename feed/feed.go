package feed

import (
	"bufio"
	"encoding/xml"
	"io"
	"net/http"
	"os"
	"strings"
)

func GetURLs(urlfile string) ([]string, error) {
	feeds := make([]string, 0)
	f, err := os.Open(urlfile)
	if err != nil {
		return feeds, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		feeds = append(feeds, strings.TrimSpace(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return feeds, err
	}
	return feeds, nil
}

func parseFeed(feed []byte) (Feed, error) {
	var feedR RSS
	var feedA Atom
	err := xml.Unmarshal(feed, &feedR)
	if err != nil {
		err = xml.Unmarshal(feed, &feedA)
		if err != nil {
			return nil, err
		}
		return feedA, nil
	}
	return feedR, nil
}

func readLocalFeed(path string) ([]byte, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}

func getRemoteFeed(url string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, nil
	}
	resp, err := client.Do(req)
	if err != nil {
		return []byte{}, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func processFeed(url string) (Feed, error) {
	var f Feed
	b, err := getRemoteFeed(url)
	if err != nil {
		return f, err
	}
	feed, err := parseFeed(b)
	if err != nil {
		return f, err
	}
	return feed, nil
}

func GetFeeds(urls []string, c chan Feed) error {
	for _, url := range urls {
		go func() {
			f, err := processFeed(url)
			if err != nil {
				return err
			}
			c <- feed
		}()
	}
}
