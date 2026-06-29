package sessions

import (
	"github.com/redis/go-redis/v9"
)

type Module struct {
	service *Service
}

func NewModule(rdb *redis.Client) *Module {
	service := NewService(rdb)
	return &Module{
		service: service,
	}
}

func (m *Module) Service() *Service {
	return m.service
}
