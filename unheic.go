package unheic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrBadRequest = errors.New("bad request")

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type ClientOption func(*Client)

func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) { c.httpClient = client }
}

func WithBaseURL(url string) ClientOption {
	return func(c *Client) { c.baseURL = url }
}

func NewClient(opts ...ClientOption) *Client {
	c := &Client{
		httpClient: &http.Client{},
		baseURL:    "http://localhost:8080",
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Convert(ctx context.Context, src io.Reader) (io.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/convert", src)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "image/heic")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusBadRequest {
		return nil, ErrBadRequest
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return resp.Body, nil
}
