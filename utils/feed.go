package utils

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Item struct {
	Title       string `xml:"title,omitempty"`
	Link        string `xml:"link,omitempty"`
	Guid        string `xml:"guid,omitempty"`
	Description string `xml:"description,omitempty"`
	PubDate     string `xml:"pubDate,omitempty"`
}

type Channel struct {
	Title         string `xml:"title,omitempty"`
	Link          string `xml:"link,omitempty"`
	Description   string `xml:"description,omitempty"`
	Language      string `xml:"language,omitempty"`
	LastBuildDate string `xml:"lastBuildDate,omitempty"`
	Item          []Item `xml:"item,omitempty"`
}

type XMLFeed struct {
	Channel Channel `xml:"channel,omitempty"`
}

func FetchFeed(url *url.URL) (*XMLFeed, error) {
	resp, err := http.Get(url.String())
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("expected status=%d, got=%d", resp.StatusCode, http.StatusOK))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	xmlFeed := &XMLFeed{}
	err = xml.Unmarshal(body, xmlFeed)
	if err != nil {
		return nil, err
	}

	return xmlFeed, nil
}
