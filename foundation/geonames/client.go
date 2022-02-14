package geonames

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/pkg/errors"
)

const (
	apiEndpoint = "http://api.geonames.org/searchJSON"
)

// Client is an API client for geonames.org
type Client struct {
	username string
}

// NewClient create a new client using username
func NewClient(username string) *Client {
	return &Client{username: username}
}

// FindOneLocation returns first location from the results
func (c *Client) FindOneLocation(ctx context.Context, name string) (*Geoname, error) {
	uv := url.Values{}
	uv.Add("q", name)
	uv.Add("maxRows", "1")
	uv.Add("username", c.username)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", apiEndpoint, uv.Encode()), nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	var result LocationResult

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Geonames) == 0 {
		return nil, errors.New("not found")
	}
	return &result.Geonames[0], nil
}
