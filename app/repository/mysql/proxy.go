package mysql

import (
	"context"

	"github.com/hiroyaonoe/bcop-proxy-controller/app/entity"
	"github.com/hiroyaonoe/bcop-proxy-controller/app/repository"
	"github.com/jmoiron/sqlx"
)

type Proxy struct {
	db *sqlx.DB
}

var _ repository.Proxy = &Proxy{}

func NewProxy(db *sqlx.DB) *Proxy {
	return &Proxy{db: db}
}

func (p *Proxy) Get(ctx context.Context, proxyID string) (entity.Proxy, error) {
	return entity.Proxy{}, nil
}

func (p *Proxy) Upsert(ctx context.Context, proxy entity.Proxy) error {
	var query string
	// Endpointがnullの場合はupdateのみを実行
	if len(proxy.Endpoint) == 0 {
		query = `
			UPDATE proxies
			SET
				activate = :activate
			WHERE
				proxy_id = :proxy_id
		`
	} else {
		query = `
			INSERT INTO proxies (proxy_id, endpoint, activate)
			VALUES (:proxy_id, :endpoint, :activate)
			ON DUPLICATE KEY
			UPDATE
				endpoint = VALUES(endpoint),
				activate = VALUES(activate)
		`
	}
	_, err := p.db.NamedExecContext(ctx, query, &proxy)
	return err
}

func (p *Proxy) Delete(ctx context.Context, proxyID string) error {
	query := `
		DELETE FROM proxies
		WHERE proxy_id = ?
	`
	_, err := p.db.ExecContext(ctx, query, proxyID)
	return err
}
