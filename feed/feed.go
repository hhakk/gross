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

func parseFeed(b []byte) (Feed, error) {
	var feedR RSS
	err := xml.Unmarshal(b, &feedR)
	if err == nil {
		return feedR, nil
	}
	var feedA Atom
	err = xml.Unmarshal(b, &feedA)
	if err == nil {
		return feedA, nil
	}
	return nil, err
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

func processFeed(url string, index int, c chan FeedMessage) {
	b, err := getRemoteFeed(url)
	if err != nil {
		c <- FeedMessage{Feed: nil, Error: err, Index: index}
		return
	}
	f, err := parseFeed(b)
	if err != nil {
		c <- FeedMessage{Feed: nil, Error: err, Index: index}
		return
	}
	c <- FeedMessage{Feed: &f, Error: nil, Index: index}
	return
}

type FeedMessage struct {
	Feed  *Feed
	Error error
	Index int
}

func GetFeeds(urls []string, c chan FeedMessage) {
	for i, url := range urls {
		go processFeed(url, i, c)
	}
}
