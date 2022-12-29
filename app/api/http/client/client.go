package client

import (
	"context"
	"io"
	"net/http"
	"net/url"
)

type Client struct {
	client *http.Client
	base   string
}

func NewClient(client *http.Client, base string) *Client {
	return &Client{
		client: client,
		base:   base,
	}
}

func (c *Client) Get(ctx context.Context, pathElem []string) (*http.Response, error) {
	reqURL, err := url.JoinPath(c.base, pathElem...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *Client) Put(ctx context.Context, pathElem []string, contentType string, body io.Reader) (*http.Response, error) {
	reqURL, err := url.JoinPath(c.base, pathElem...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, reqURL, body)
	if err != nil {
		return nil, err
	}
	if len(contentType) != 0 {
		req.Header.Set("Content-Type", contentType)
	}
	return c.client.Do(req)
}

func (c *Client) Delete(ctx context.Context, pathElem []string) (*http.Response, error) {
	reqURL, err := url.JoinPath(c.base, pathElem...)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, reqURL, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}
