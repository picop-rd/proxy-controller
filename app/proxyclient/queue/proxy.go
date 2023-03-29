package queue

import (
	"context"

	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy-controller/app/proxyclient"
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
		return nil
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
