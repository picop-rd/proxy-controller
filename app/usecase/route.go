package usecase

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/repository"
)

type Route struct {
	repo repository.Repository
}

func NewRoute(repo repository.Repository) *Route {
	return &Route{
		repo: repo,
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
	// TODO: キューにrouteを詰める
	return nil
}

func (r *Route) Delete(ctx context.Context, routes []entity.Route) error {
	err := r.repo.Route.Delete(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to delete routes from repository: %w", err)
	}
	// TODO: キューにrouteを詰める
	return nil
}
