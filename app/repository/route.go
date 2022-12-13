package repository

import (
	"context"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
)

type Route interface {
	GetWithProxyID(ctx context.Context, proxyID string) ([]entity.Route, error)
	Upsert(ctx context.Context, routes []entity.Route) error
	Delete(ctx context.Context, routes []entity.Route) error
}
