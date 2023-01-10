package usecase

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/proxyclient"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/repository"
)

type Proxy struct {
	repo   repository.Repository
	client proxyclient.Client
}

func NewProxy(repo repository.Repository, client proxyclient.Client) *Proxy {
	return &Proxy{
		repo:   repo,
		client: client,
	}
}

func (p *Proxy) Register(ctx context.Context, proxy entity.Proxy) error {
	proxy.Activate = false
	if err := proxy.Validate(); err != nil {
		return fmt.Errorf("invalid proxy: %w", err)
	}

	err := p.repo.Proxy.Upsert(ctx, proxy)
	if err != nil {
		return fmt.Errorf("failed to register proxy to repository: %w", err)
	}

	err = p.client.Proxy.Deactivate(ctx, proxy.ProxyID)
	if err != nil {
		return fmt.Errorf("failed to deactivate proxy by proxyclient: %w", err)
	}
	return nil
}

func (p *Proxy) Activate(ctx context.Context, proxyID string) error {
	proxy := entity.Proxy{
		ProxyID:  proxyID,
		Endpoint: "",
		Activate: true,
	}
	err := p.repo.Proxy.Upsert(ctx, proxy)
	if err != nil {
		return fmt.Errorf("failed to activate proxy on repository: %w", err)
	}

	proxy, err = p.repo.Proxy.Get(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to get proxy from repository: %w", err)
	}
	err = p.client.Proxy.Activate(ctx, proxy)
	if err != nil {
		return fmt.Errorf("failed to activate proxy by proxyclient: %w", err)
	}

	routes, err := p.repo.Route.GetWithProxyID(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to get routes from repository: %w", err)
	}
	if len(routes) == 0 {
		return nil
	}
	err = p.client.Route.Register(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to register routes by proxyclient: %w", err)
	}

	return nil
}

func (p *Proxy) Delete(ctx context.Context, proxyID string) error {
	err := p.repo.Proxy.Delete(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to delete proxy from repository: %w", err)
	}

	err = p.client.Proxy.Deactivate(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to deactivate proxy by proxyclient: %w", err)
	}
	return nil
}
