package usecase

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/repository"
)

type Route struct {
	route repository.Route
}

func NewRoute(route repository.Route) *Route {
	return &Route{route: route}
}

func (r *Route) Register(ctx context.Context, routes []entity.Route) error {
	for _, route := range routes {
		err := route.Validate()
		if err != nil {
			return fmt.Errorf("invalid route: %w", err)
		}
	}
	err := r.route.Upsert(ctx, routes)
	if err != nil {
		return fmt.Errorf("failed to register routes to registory: %w", err)
	}
	// TODO: キューにrouteを詰める
	return nil
}
