package entity

type Route struct {
	ProxyID     string `db:"proxy_id" json:"proxy_id"`
	EnvID       string `db:"env_id" json:"env_id"`
	Destination string `db:"destination" json:"destination"`
}
