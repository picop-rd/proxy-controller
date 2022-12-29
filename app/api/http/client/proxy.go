package client

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
)

type Proxy struct {
	client *Client
}

func NewProxy(client *Client) *Proxy {
	return &Proxy{
		client: client,
	}
}

func (p *Proxy) Register(ctx context.Context, proxy entity.Proxy) error {
	body := fmt.Sprintf("{\"endpoint\": \"%s\"}", proxy.Endpoint)
	resp, err := p.client.Put(ctx, []string{"proxy", proxy.ProxyID, "register"}, "application/json", strings.NewReader(body))
	if err != nil {
		return fmt.Errorf("proxy controller client: Register: failed to request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxy controller client: Register: response status code is not 200, status: %s", resp.Status)
	}
	return nil
}

func (p *Proxy) Activate(ctx context.Context, proxyID string) error {
	resp, err := p.client.Put(ctx, []string{"proxy", proxyID, "activate"}, "application/json", nil)
	if err != nil {
		return fmt.Errorf("proxy controller client: Activate: failed to request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxy controller client: Activate: response status code is not 200, status: %s", resp.Status)
	}
	return nil
}

func (p *Proxy) Delete(ctx context.Context, proxyID string) error {
	resp, err := p.client.Delete(ctx, []string{"proxy", proxyID})
	if err != nil {
		return fmt.Errorf("proxy controller client: Delete: failed to request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("proxy controller client: Delete: response status code is not 200, status: %s", resp.Status)
	}
	return nil
}
