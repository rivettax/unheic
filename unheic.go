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
}

var ErrBadRequest = errors.New("bad request")

func NewClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

func (c *Client) Convert(ctx context.Context, src io.Reader) (io.Reader, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", "http://localhost:8080/convert", src)
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
