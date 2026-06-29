package auth

import (
	"github.com/Rowkash/go-gin-auth/internal/config"
	"github.com/Rowkash/go-gin-auth/internal/sessions"
	"github.com/Rowkash/go-gin-auth/internal/users"
	"github.com/gin-gonic/gin"
)

type ModuleDeps struct {
	JwtCfg          config.JWTConfig
	UsersService    *users.Service
	SessionsService *sessions.Service
}

type Module struct {
	service         *Service
	handler         *Handler
	tokenService    *TokenService
	sessionsService *sessions.Service
}

func NewModule(deps ModuleDeps) *Module {
	tokenService := NewTokenService(deps.JwtCfg.Secret, deps.JwtCfg.ExpiresIn)
	service := NewService(deps.UsersService, tokenService, deps.SessionsService)

	return &Module{
		service:         service,
		handler:         NewHandler(service),
		tokenService:    tokenService,
		sessionsService: deps.SessionsService,
	}
}

func (m *Module) Middleware() gin.HandlerFunc {
	return middleware(m.tokenService)
}
