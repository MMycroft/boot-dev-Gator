// Package feed holds feed stuff
package feed

import (
	"context"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string `xml:"title"`
		Link        string `xml:"link"`
		Description string `xml:"description"`
		Item        []struct {
			Title       string `xml:"title"`
			Link        string `xml:"link"`
			Description string `xml:"description"`
			PubDate     string `xml:"pubDate"`
		} `xml:"item"`
	} `xml:"channel"`
}

func FetchFeed(ctx context.Context, url string) (*RSSFeed, error) {
	rssFeed := &RSSFeed{}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return rssFeed, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("User-Agent", "gator")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rssFeed, fmt.Errorf("error sending request: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return rssFeed, fmt.Errorf("error reading response body: %w", err)
	}

	err = xml.Unmarshal(b, rssFeed)
	if err != nil {
		return rssFeed, fmt.Errorf("error unmarshaling response bytes: %w", err)
	}

	rssFeed.Unescape()

	return rssFeed, nil
}

func (rss *RSSFeed) Unescape() {
	rss.Channel.Title = html.UnescapeString(rss.Channel.Title)
	rss.Channel.Link = html.UnescapeString(rss.Channel.Link)
	rss.Channel.Description = html.UnescapeString(rss.Channel.Description)
	for _, item := range rss.Channel.Item {
		item.Title = html.UnescapeString(item.Title)
		item.Link = html.UnescapeString(item.Link)
		item.Description = html.UnescapeString(item.Description)
	}
}
