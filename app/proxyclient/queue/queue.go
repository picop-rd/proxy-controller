package queue

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy-controller/app/proxyclient"
	"github.com/picop-rd/proxy/app/admin/api/http/client"
	proxyEntity "github.com/picop-rd/proxy/app/entity"
	"github.com/rs/zerolog/log"
)

var (
	ErrContextCanceled = errors.New("proxy client queue: context canceled")
)

type Queue struct {
	queue      *Map[*item]
	client     *http.Client
	interval   time.Duration
	ctx        context.Context
	cancelFunc context.CancelFunc
	quit       chan struct{}
}

func NewQueue(client *http.Client, interval time.Duration) *Queue {
	if client == nil {
		client = http.DefaultClient
	}
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	return &Queue{
		queue:      NewMap[*item](),
		client:     client,
		interval:   interval,
		ctx:        ctx,
		cancelFunc: cancel,
		quit:       make(chan struct{}, 2),
	}
}

func (q *Queue) Start() error {
	return q.start()
}

func (q *Queue) Close() {
	log.Info().Msg("closing proxy client queue")
	q.cancelFunc()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	select {
	case <-q.quit:
	case <-ctx.Done():
		log.Fatal().Msg("failed to close queue")
		return
	}
}

func NewClient(queue *Queue) proxyclient.Client {
	return proxyclient.Client{
		Proxy: NewProxy(queue),
		Route: NewRoute(queue),
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
	interval := time.NewTicker(q.interval)
	log.Info().Msg("starting proxy client queue process")
	for {
		select {
		case <-q.ctx.Done():
			q.quit <- struct{}{}
			return nil
		case <-interval.C:
		}

		q.queue.Range(func(proxyID string, it *item) bool {
			log.Info().Str("proxyID", proxyID).Msg("processing proxy")
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
	it.registers.Range(func(envID string, route entity.Route) bool {
		log.Debug().Str("envID", envID).Msg("processing to register route")
		env := proxyEntity.Env{
			EnvID:       route.EnvID,
			Destination: route.Destination,
		}
		envs = append(envs, env)
		return true
	})
	if len(envs) == 0 {
		return nil
	}
	err := it.envCli.Register(ctx, envs)
	if err != nil {
		return err
	}
	for _, e := range envs {
		it.registers.Del(e.EnvID)
	}
	return nil
}

func processDeletes(ctx context.Context, it *item) error {
	// sync.Mapはlenが不明なので0としておく
	envIDs := make([]string, 0)
	it.deletes.Range(func(envID string, _ entity.Route) bool {
		log.Debug().Str("envID", envID).Msg("processing to delete route")
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
