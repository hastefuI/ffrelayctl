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
	DefaultBaseURL  = "https://relay.firefox.com"
	APIBasePath     = "/api/v1/"
	DefaultTimeout  = 30 * time.Second
	ContentTypeJson = "application/json"
)

type Client struct {
	BaseURL    string
	Token      string
	UserAgent  string
	HTTPClient *http.Client
	ctx        context.Context
}

type ClientOption func(*Client)

func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = strings.TrimSuffix(url, "/")
	}
}

func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.HTTPClient = httpClient
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.HTTPClient.Timeout = timeout
	}
}

func WithUserAgent(userAgent string) ClientOption {
	return func(c *Client) {
		c.UserAgent = userAgent
	}
}

func WithContext(ctx context.Context) ClientOption {
	return func(c *Client) {
		c.ctx = ctx
	}
}

func NewClient(token string, opts ...ClientOption) *Client {
	c := &Client{
		BaseURL: DefaultBaseURL,
		Token:   token,
		HTTPClient: &http.Client{
			Timeout: DefaultTimeout,
		},
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	return c.NewRequestWithContext(c.ctx, method, path, body)
}

func (c *Client) NewRequestWithContext(ctx context.Context, method, path string, body io.Reader) (*http.Request, error) {
	url := c.BaseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Token "+c.Token)
	req.Header.Set("Content-Type", ContentTypeJson)
	req.Header.Set("Accept", ContentTypeJson)
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

func (c *Client) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	return resp, nil
}

func (c *Client) Get(path string) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Post(path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPost, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Put(path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPut, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Patch(path string, body io.Reader) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodPatch, path, body)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}

func (c *Client) Delete(path string) (*http.Response, error) {
	req, err := c.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req)
}
