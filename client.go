package hn

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// DefaultClient is the default Hacker News API client
var DefaultClient = NewClient()

const (
	libraryVersion = "0.0.2"
	userAgent      = "hn.go/" + libraryVersion
)

// A Client communicates with the Hacker News API.
type Client struct {
	Items ItemsService
	Users UsersService
	Live  LiveService

	// BaseURL is the base url for Hacker News API.
	BaseURL *url.URL

	// User agent used for HTTP requests to Hacker News API.
	UserAgent string

	// HTTP client used to communicate with the Hacker News API.
	httpClient *http.Client
}

// NewClient returns a new Hacker News API client.
// If no *http.Client were provided then http.DefaultClient is used.
func NewClient(httpClients ...*http.Client) *Client {
	var httpClient *http.Client

	if len(httpClients) > 0 && httpClients[0] != nil {
		httpClient = httpClients[0]
	} else {
		cloned := *http.DefaultClient
		httpClient = &cloned
	}

	c := &Client{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "hacker-news.firebaseio.com",
			Path:   "/v0/",
		},
		UserAgent:  userAgent,
		httpClient: httpClient,
	}

	c.Items = &itemsService{c}
	c.Users = &usersService{c}
	c.Live = &liveService{c}

	return c
}

// NewRequest creates an API request.
func (c *Client) NewRequest(ctx context.Context, s string) (*http.Request, error) {
	rel, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// decoded and stored in the value pointed to by v, or returned as an error if
// an API error has occurred.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	// Make sure to close the connection after replying to this request
	req.Close = true

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}

	if err != nil {
		return nil, fmt.Errorf("error reading response from %s %s: %s", req.Method, req.URL.RequestURI(), err)
	}

	return resp, nil
}
