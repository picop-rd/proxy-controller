package queue

import (
	"github.com/picop-rd/proxy-controller/app/entity"
	"github.com/picop-rd/proxy/app/admin/api/http/client"
)

type item struct {
	envCli    *client.Env
	registers *Map[entity.Route]
	deletes   *Map[entity.Route]
}

func newItem(cli *client.Client) *item {
	return &item{
		envCli:    client.NewEnv(cli),
		registers: NewMap[entity.Route](),
		deletes:   NewMap[entity.Route](),
	}
}

func (i *item) Register(route entity.Route) {
	envID := route.EnvID
	i.deletes.Del(envID)
	i.registers.Set(envID, route)
}

func (i *item) Delete(route entity.Route) {
	envID := route.EnvID
	i.registers.Del(envID)
	i.deletes.Set(envID, route)
}
