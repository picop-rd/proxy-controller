package usecase

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/repository"
)

type Proxy struct {
	repo repository.Repository
}

func NewProxy(repo repository.Repository) *Proxy {
	return &Proxy{
		repo: repo,
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

	_, err = p.repo.Route.GetWithProxyID(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to get routes from repository: %w", err)
	}
	// TODO: proxyclientを通したリクエストによってrouteを追加する(キューにつめる)
	return nil
}

func (p *Proxy) Delete(ctx context.Context, proxyID string) error {
	err := p.repo.Proxy.Delete(ctx, proxyID)
	if err != nil {
		return fmt.Errorf("failed to delete proxy from repository: %w", err)
	}
	// TODO: proxyclientへのキューから該当proxyのrouteを削除する
	return nil
}
