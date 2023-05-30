package queue

import (
	"context"
	"fmt"

	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy-controller/app/proxyclient"
)

type Route struct {
	queue *Queue
}

var _ proxyclient.Route = &Route{}

func NewRoute(queue *Queue) *Route {
	return &Route{
		queue: queue,
	}
}

func (r *Route) Register(ctx context.Context, routes []entity.Route) error {
	select {
	case <-ctx.Done():
		return ErrContextCanceled
	default:
	}

	for _, route := range routes {
		proxyID := route.ProxyID
		it, ok := r.queue.get(proxyID)
		if !ok {
			// You can also register the route before activating the proxy.
			continue
		}
		it.Register(route)
	}
	return nil
}

func (r *Route) Delete(ctx context.Context, routes []entity.Route) error {
	select {
	case <-ctx.Done():
		return ErrContextCanceled
	default:
	}

	for _, route := range routes {
		proxyID := route.ProxyID
		it, ok := r.queue.get(proxyID)
		if !ok {
			return fmt.Errorf("proxy client queue: not found item in queue, proxyID: %s", proxyID)
		}
		it.Delete(route)
	}
	return nil
}
