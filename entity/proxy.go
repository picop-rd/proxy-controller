package entity

type Proxy struct {
	ProxyID  string `db:"proxy_id" json:"-"`
	Endpoint string `db:"endpoint" json:"endpoint"`
	Activate bool   `db:"activate" json:"-"`
}

func (p Proxy) Validate() error {
	if len(p.ProxyID) == 0 || len(p.Endpoint) == 0 {
		return ErrInvalid
	}
	return nil
}
