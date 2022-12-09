package entity

import "github.com/rs/zerolog"

type Route struct {
	ProxyID     string `db:"proxy_id" json:"proxy_id"`
	EnvID       string `db:"env_id" json:"env_id"`
	Destination string `db:"destination" json:"destination"`
}

func (r Route) MarshalZerologObject(e *zerolog.Event) {
	e.Str("ProxyID", r.ProxyID).
		Str("EnvID", r.EnvID).
		Str("Destination", r.Destination)
}
