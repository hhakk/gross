package feed

import (
	"bytes"
	"bufio"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/net/html/charset"
	"github.com/spf13/viper"
)

type FeedSpec struct {
	URL string
	Cmd string
	AltName string
}

func GetURLs(urlfile string) ([]FeedSpec, error) {
	feeds := make([]FeedSpec, 0)
	f, err := os.Open(urlfile)
	if err != nil {
		return feeds, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		base := strings.SplitN(strings.TrimSpace(scanner.Text()), " ", 2)
		if base[0] == "" { continue }
		fs := FeedSpec{}
		for i, arg := range base {
			if i == 0 {
				// handle cmd like filter:cmd:url 
				_, cmdurl, ok := strings.Cut(arg, "filter:")
				if ok {
					cmd, url, ok := strings.Cut(cmdurl, ":")
					if ok {
						fs.Cmd = cmd
						fs.URL = url
						continue
					}
				}
				fs.URL = arg
				continue
			}
			fs.AltName = strings.ReplaceAll(arg, "\"", "")
		}
		feeds = append(feeds, fs)
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

func getOldFeed(url FeedSpec) (Feed, error) {
	// check previous feed
	oldfp := filepath.Join(viper.GetString("cachedir"), fmt.Sprintf("%x", md5.Sum([]byte(url.URL))))
	oldb, err := os.ReadFile(oldfp)
	if err != nil {
		return nil, err
	}
	return parseFeed(oldb, url)
}

func parseFeed(b []byte, url FeedSpec) (Feed, error) {

	baseurl, err := neturl.Parse(url.URL)
	if err != nil {
		return nil, err
	}
	if url.Cmd != "" {
		cmd := exec.Command(url.Cmd)
		stdin, err := cmd.StdinPipe()
		if err != nil {
			return nil, err
		}
		go func() {
			defer stdin.Close()
			io.WriteString(stdin, string(b))
		}()
		out, err := cmd.Output()
		if err != nil {
			return nil, err
		}
		b = []byte(out)
	}
	var feedR RSS
	reader := bytes.NewReader(b)
	decoder := xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&feedR)
	if err == nil {
		feedR.url = url.URL
		for i := range feedR.XChannel.XItems {
			feedR.XChannel.XItems[i].url = fmt.Sprintf("%s://%s", baseurl.Scheme, baseurl.Host)
		}
		if url.AltName != "" {
			feedR.AltTitle = url.AltName
		}
		return &feedR, nil
	}
	var feedA Atom
	reader = bytes.NewReader(b)
	decoder = xml.NewDecoder(reader)
	decoder.CharsetReader = charset.NewReaderLabel
	err2 := decoder.Decode(&feedA)
	if err2 == nil {
		feedA.url = url.URL
		for i := range feedA.XEntries {
			feedA.XEntries[i].url = fmt.Sprintf("%s://%s", baseurl.Scheme, baseurl.Host)
		}
		if url.AltName != "" {
			feedA.AltTitle = url.AltName
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

func processFeed(url FeedSpec, index int, c chan FeedMessage) {
	b, err := getRemoteFeed(url.URL)
	if err != nil {
		c <- FeedMessage{Feed: nil, Error: err, Index: index, URL: url.URL}
		return
	}
	f, err := parseFeed(b, url)
	if err != nil {
		c <- FeedMessage{Feed: nil, Error: err, Index: index, URL: url.URL}
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
	c <- FeedMessage{Feed: &f, Error: nil, Index: index, URL: url.URL}
	return
}

type FeedMessage struct {
	Feed  *Feed
	Error error
	Index int
	URL string
}

func GetFeeds(urls []FeedSpec, c chan FeedMessage) {
	for i, url := range urls {
		go processFeed(url, i, c)
	}
}
