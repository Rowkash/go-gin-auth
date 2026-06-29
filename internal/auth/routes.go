package auth

import (
	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(router *gin.RouterGroup) {

	public := router.Group("/auth")
	{
		public.POST("/registration", m.handler.register)
		public.POST("/login", m.handler.login)
	}

	authorized := router.Group("/auth")
	authorized.Use(m.Middleware())
	{
		authorized.DELETE("/logout", m.handler.logout)

	}
}
