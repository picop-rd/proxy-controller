package queue

import (
	"context"
	"fmt"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/proxyclient"
)

type Proxy struct {
	queue *Queue
}

var _ proxyclient.Proxy = &Proxy{}

func NewProxy(queue *Queue) *Proxy {
	return &Proxy{
		queue: queue,
	}
}

func (p *Proxy) Activate(ctx context.Context, proxy entity.Proxy) error {
	select {
	case <-ctx.Done():
		return ErrContextCanceled
	default:
	}

	proxyID := proxy.ProxyID
	if _, ok := p.queue.get(proxyID); ok {
		return fmt.Errorf("proxy client queue: cannot activate existing proxy, proxyID: %s", proxyID)
	}
	p.queue.add(proxyID, proxy.Endpoint)
	return nil
}

func (p *Proxy) Deactivate(ctx context.Context, proxyID string) error {
	select {
	case <-ctx.Done():
		return ErrContextCanceled
	default:
	}

	p.queue.del(proxyID)
	return nil
}
