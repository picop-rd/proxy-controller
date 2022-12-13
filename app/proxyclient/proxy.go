package proxyclient

import (
	"context"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
)

type Proxy interface {
	Activate(ctx context.Context, proxy entity.Proxy) error
	Deactivate(ctx context.Context, proxyID string) error
}
