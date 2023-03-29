package usecase

import (
	"context"
	"fmt"

	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy-controller/app/proxyclient"
	"github.com/picop-rd/proxy-controller/app/repository"
)

type Route struct {
	repo   repository.Repository
	client proxyclient.Client
}

func NewRoute(repo repository.Repository, client proxyclient.Client) *Route {
	return &Route{
		repo:   repo,
		client: client,
	}
}

func (r *Route) Register(ctx context.Context, routes []entity.Route) error {
	for _, route := range routes {
		err := route.Validate()
		if err != nil {
			return fmt.Errorf("invalid route: %w", err)
		}
	}
	err := r.repo.Route.Upsert(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to register routes to repository: %w", err)
	}
	err = r.client.Route.Register(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to register routes by proxyclient: %w", err)
	}
	return nil
}

func (r *Route) Delete(ctx context.Context, routes []entity.Route) error {
	err := r.repo.Route.Delete(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to delete routes from repository: %w", err)
	}
	err = r.client.Route.Delete(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to delete routes by proxyclient: %w", err)
	}
	return nil
}
