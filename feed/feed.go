package feed

import (
	"bufio"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
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

func SaveFeed(feed Feed) error {
	url := feed.URL()
	fp := filepath.Join(viper.GetString("cachedir"), fmt.Sprintf("%x", md5.Sum([]byte(url))))
	switch feed.(type) {
	case *RSS:
		b, err := xml.Marshal(&feed)
		if err != nil {
			return err
		}
		err = os.WriteFile(fp, b, 0666)
		if err != nil {
			return err
		}
	case *Atom:
		b, err := xml.Marshal(&feed)
		if err != nil {
			return err
		}
		err = os.WriteFile(fp, b, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func getOldFeed(url string) (Feed, error) {
	// check previous feed
	oldfp := filepath.Join(viper.GetString("cachedir"), fmt.Sprintf("%x", md5.Sum([]byte(url))))
	oldb, err := os.ReadFile(oldfp)
	if err != nil {
		return nil, err
	}
	return parseFeed(oldb, url)
}

func parseFeed(b []byte, url string) (Feed, error) {

	baseurl, err := neturl.Parse(url)
	if err != nil {
		return nil, err
	}
	var feedR RSS
	err = xml.Unmarshal(b, &feedR)
	if err == nil {
		feedR.url = url
		for i := range feedR.XChannel.XItems {
			feedR.XChannel.XItems[i].url = fmt.Sprintf("%s://%s", baseurl.Scheme, baseurl.Host)
		}
		return &feedR, nil
	}
	var feedA Atom
	err = xml.Unmarshal(b, &feedA)
	if err == nil {
		feedA.url = url
		for i := range feedA.XEntries {
			feedA.XEntries[i].url = fmt.Sprintf("%s://%s", baseurl.Scheme, baseurl.Host)
		}
		return &feedA, nil
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
	f, err := parseFeed(b, url)
	if err != nil {
		c <- FeedMessage{Feed: nil, Error: err, Index: index}
		return
	}
	oldf, err := getOldFeed(url)
	if err == nil {
		for _, old := range oldf.Items() {
			for _, cur := range f.Items() {
				if old.Title() == cur.Title() &&
					old.Description() == cur.Description() &&
					old.Link() == cur.Link() &&
					old.Content() == cur.Content() &&
					old.IsRead() {
					cur.SetRead(true)
				}
			}
		}

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
