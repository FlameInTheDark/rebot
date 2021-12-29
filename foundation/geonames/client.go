package geonames

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"net/http"
)

type Client struct {
	username string
}

func NewClient(username string) *Client {
	return &Client{username: username}
}

func (c *Client) FindOneLocation(name string) (*Geoname, error) {
	resp, err := http.Get(fmt.Sprintf("http://api.geonames.org/searchJSON?q=%v&maxRows=1&username=%v", name, c.username))
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
