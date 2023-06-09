package proxyclient

import (
	"context"

	"github.com/picop-rd/proxy-controller/app/entity"
)

type Route interface {
	Register(ctx context.Context, routes []entity.Route) error
	Delete(ctx context.Context, route []entity.Route) error
}
