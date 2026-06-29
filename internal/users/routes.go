package users

import (
	"github.com/Rowkash/go-gin-auth/internal/common/middleware"
	"github.com/Rowkash/go-gin-auth/internal/common/types"
	"github.com/gin-gonic/gin"
)

func (m *Module) RegisterRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	usersGroup := rg.Group("/users")
	usersGroup.Use(authMiddleware)

	{
		usersGroup.GET("/", middleware.RolesAuth(types.RoleAdmin), m.handler.getPage)
		usersGroup.GET("/:id", middleware.RolesAuth(types.RoleAdmin), m.handler.findById)
	}
}
