package geonames

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

type Client struct {
	username string
}

func NewClient(username string) *Client {
	return &Client{username: username}
}

func (c *Client) FindOneLocation(ctx context.Context, name string) (*Geoname, error) {
	uv := url.Values{}
	uv.Add("q", name)
	uv.Add("maxRows", "1")
	uv.Add("username", c.username)

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://api.geonames.org/searchJSON?%s", uv.Encode()), nil)
	if err != nil {
		return nil, err
	}

	client := http.DefaultClient
	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

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
