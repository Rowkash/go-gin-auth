package users

import "gorm.io/gorm"

type Module struct {
	service *Service
	handler *Handler
}

func NewModule(db *gorm.DB) *Module {
	service := NewService(db)
	return &Module{
		service: service,
		handler: NewHandler(service),
	}
}

func (m *Module) Service() *Service {
	return m.service
}

//func (m *Module) Handler() *Handler {
//	return m.handler
//}
