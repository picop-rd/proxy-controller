package queue

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy/app/admin/api/http/client"
	proxyEntity "github.com/hiroyaonoe/bcop-proxy/app/entity"
	"github.com/rs/zerolog/log"
)

var (
	ErrContextCanceled = errors.New("proxy client queue: context canceled")
)

type Queue struct {
	queue      *Map[*item]
	client     *http.Client
	ctx        context.Context
	cancelFunc context.CancelFunc
	quit       chan struct{}
}

func NewQueue(client *http.Client) *Queue {
	if client == nil {
		client = http.DefaultClient
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &Queue{
		queue:      NewMap[*item](),
		client:     client,
		ctx:        ctx,
		cancelFunc: cancel,
		quit:       make(chan struct{}, 2),
	}
}

func (q *Queue) Start() error {
	return q.start()
}

func (q *Queue) Close() {
	q.cancelFunc()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	for i := 0; i < 2; i++ {
		select {
		case <-q.quit:
		case <-ctx.Done():
			log.Fatal().Msg("failed to close queue")
			return
		}
	}
}

func (q *Queue) get(proxyID string) (*item, bool) {
	return q.queue.Get(proxyID)
}

func (q *Queue) add(proxyID, endpoint string) {
	it := newItem(client.NewClient(q.client, endpoint))
	q.queue.Set(proxyID, it)
}

func (q *Queue) del(proxyID string) {
	q.queue.Del(proxyID)
}

func (q *Queue) start() error {
	for {
		select {
		case <-q.ctx.Done():
			q.quit <- struct{}{}
			return nil
		default:
		}

		q.queue.Range(func(proxyID string, it *item) bool {
			err := process(q.ctx, it)
			if err != nil {
				log.Error().Err(err).Str("proxyID", proxyID).Msg("failed to process queue")
			}
			return true
		})
	}
}
func process(ctx context.Context, it *item) error {
	err := processRegisters(ctx, it)
	if err != nil {
		return err
	}
	err = processDeletes(ctx, it)
	if err != nil {
		return err
	}
	return nil
}

func processRegisters(ctx context.Context, it *item) error {
	// sync.Mapはlenが不明なので0としておく
	envs := make([]proxyEntity.Env, 0)
	bufRoutes := make([]entity.Route, 0)
	it.registers.Range(func(envID string, route entity.Route) bool {
		env := proxyEntity.Env{
			EnvID:       route.EnvID,
			Destination: route.Destination,
		}
		envs = append(envs, env)
		bufRoutes = append(bufRoutes, route)
		return true
	})
	err := it.envCli.Register(ctx, envs)
	if err != nil {
		for _, r := range bufRoutes {
			it.registers.Set(r.EnvID, r)
		}
		return err
	}
	return nil
}

func processDeletes(ctx context.Context, it *item) error {
	// sync.Mapはlenが不明なので0としておく
	envIDs := make([]string, 0)
	it.deletes.Range(func(envID string, _ entity.Route) bool {
		envIDs = append(envIDs, envID)
		return true
	})
	for _, envID := range envIDs {
		err := it.envCli.Delete(ctx, envID)
		if err != nil {
			return err
		}
		it.deletes.Del(envID)
	}
	return nil
}
