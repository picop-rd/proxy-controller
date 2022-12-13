package repository

import (
	"context"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
)

type Proxy interface {
	Get(ctx context.Context, proxyID string) (entity.Proxy, error)
	Upsert(ctx context.Context, proxy entity.Proxy) error
	Delete(ctx context.Context, proxyID string) error
}
