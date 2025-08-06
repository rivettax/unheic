package unheic

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	httpClient *http.Client
	baseURL    string
}

type ClientParams struct {
	// BaseURL is the base URL of the unheic server. Default is http://localhost:8080
	BaseURL string

	// HTTPClient is the HTTP client to use. Default is http.Client{}
	HTTPClient *http.Client
}

var ErrBadRequest = errors.New("bad request")

func NewClient(params ClientParams) *Client {
	if params.HTTPClient == nil {
		params.HTTPClient = &http.Client{}
	}

	if params.BaseURL == "" {
		params.BaseURL = "http://localhost:8080"
	}

	return &Client{
		httpClient: params.HTTPClient,
		baseURL:    params.BaseURL,
	}
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
