// Package api is a client for the Firefox Relay API.
//
// Create a Client with NewClient and a Firefox Relay API token. Every call
// takes a context and returns the decoded response or an error. A response
// with a status of 400 or above is returned as an *APIError carrying the
// status code and the raw body.
package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the host used unless WithBaseURL says otherwise.
	DefaultBaseURL = "https://relay.firefox.com"
	// APIBasePath is the path prefix shared by every endpoint.
	APIBasePath = "/api/v1/"
	// DefaultTimeout is the HTTP timeout used unless WithTimeout says otherwise.
	DefaultTimeout = 30 * time.Second
	// ContentTypeJSON is the media type sent and accepted on every request.
	ContentTypeJSON = "application/json"
)

// Client holds the token and the HTTP settings used for API calls.
// Use NewClient to create one.
type Client struct {
	BaseURL    string
	Token      string
	UserAgent  string
	HTTPClient *http.Client
}

// ClientOption configures a Client in NewClient.
type ClientOption func(*Client)

// WithBaseURL sets the API host, dropping any trailing slash. It is meant for
// tests and for pointing at a host other than Firefox Relay.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(url, "/")
	}
}

// WithHTTPClient replaces the HTTP client used for every request. It discards
// a timeout set by an earlier WithTimeout, so pass WithTimeout after it.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

// WithTimeout sets the timeout on the current HTTP client. WithHTTPClient
// replaces that client, so pass this option after it.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

// WithUserAgent sets the User-Agent header. The header is left off when the
// value is empty.
func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

// NewClient returns a Client that authenticates with token. Options are
// applied in order over DefaultBaseURL and DefaultTimeout.
func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// NewRequest builds a request for path with the authorization, content type
// and accept headers set. path is appended to the base URL as given, so it
// should start with APIBasePath.
func (c *Client) NewRequest(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.Token)
	req.Header.Set("Content-Type", ContentTypeJSON)
	req.Header.Set("Accept", ContentTypeJSON)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends req and returns the response. It does not look at the status code,
// so an error response comes back with a nil error.
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

// Get sends a GET request to path. The caller closes the response body.
func (c *Client) Get(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Post sends a POST request to path. The caller closes the response body.
func (c *Client) Post(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Put sends a PUT request to path. The caller closes the response body.
func (c *Client) Put(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Patch sends a PATCH request to path. The caller closes the response body.
func (c *Client) Patch(ctx context.Context, path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

// Delete sends a DELETE request to path. The caller closes the response body.
func (c *Client) Delete(ctx context.Context, path string) (*http.Response, error) {
	req, err := c.NewRequest(ctx, http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
