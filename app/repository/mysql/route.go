package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy-controller/app/repository"
)

type Route struct {
	db *sqlx.DB
}

var _ repository.Route = &Route{}

func NewRoute(db *sqlx.DB) *Route {
	return &Route{db: db}
}

func (r *Route) GetWithProxyID(ctx context.Context, proxyID string) ([]entity.Route, error) {
	var routes []entity.Route
	query := `
		SELECT
			proxy_id,
			env_id,
			destination
		FROM routes
		WHERE proxy_id = ?
	`
	err := r.db.SelectContext(ctx, &routes, query, proxyID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return []entity.Route{}, nil
		}
		return nil, err
	}
	return routes, nil
}

func (r *Route) Upsert(ctx context.Context, routes []entity.Route) error {
	query := `
		INSERT INTO routes (proxy_id, env_id, destination)
		VALUES (:proxy_id, :env_id, :destination)
		ON DUPLICATE KEY
		UPDATE
			destination = VALUES(destination)
	`
	_, err := r.db.NamedExecContext(ctx, query, routes)
	return err
}

func (r *Route) Delete(ctx context.Context, routes []entity.Route) error {
	// TODO: Bulk Delete
	query := `
		DELETE FROM routes
		WHERE 
			proxy_id = :proxy_id
		AND
			env_id = :env_id
	`
	for _, route := range routes {
		_, err := r.db.NamedExecContext(ctx, query, route)
		if err != nil {
			return err
		}
	}
	return nil
}
