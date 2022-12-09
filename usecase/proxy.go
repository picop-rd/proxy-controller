package usecase

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/repository"
)

type Proxy struct {
	proxy repository.Proxy
	route repository.Route
}

func NewProxy(proxy repository.Proxy, route repository.Route) *Proxy {
	return &Proxy{
		proxy: proxy,
		route: route,
	}
}

func (p *Proxy) Register(ctx context.Context, proxy entity.Proxy) error {
	proxy.Activate = false
	if err := proxy.Validate(); err != nil {
		return fmt.Errorf("invalid proxy: %w", err)
	}

	err := p.proxy.Upsert(ctx, proxy)
	if err != nil {
		return fmt.Errorf("failed to register proxy to repository: %w", err)
	}
	return nil
}

func (p *Proxy) Activate(ctx context.Context, proxyID string) ([]entity.Route, error) {
	proxy := entity.Proxy{
		ProxyID:  proxyID,
		Endpoint: "",
		Activate: true,
	}
	err := p.proxy.Upsert(ctx, proxy)
	if err != nil {
		return nil, fmt.Errorf("failed to activate proxy on repository: %w", err)
	}

	routes, err := p.route.GetWithProxyID(ctx, proxyID)
	if err != nil {
		return nil, fmt.Errorf("failed to get routes from repository: %w", err)
	}
	// TODO: proxyclientを通したリクエストによってrouteを追加する(キューにつめる)
	return routes, nil
}

func (p *Proxy) Delete(ctx context.Context, proxyID string) error {
	return nil
}
