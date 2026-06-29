package middleware

import (
	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/Rowkash/go-gin-auth/internal/common/types"
	"github.com/gin-gonic/gin"
)

func RolesAuth(allowedRoles ...types.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		userVal, exists := c.Get(types.UserKey)
		if !exists {
			_ = c.Error(common.NewUnauthorizedError("Access denied: Unauthenticated"))
			c.Abort()
			return
		}

		user := userVal.(types.UserContext)

		currentRole := types.UserRole(user.Role)

		roleAllowed := false
		for _, role := range allowedRoles {
			if currentRole == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			_ = c.Error(common.NewForbiddenError("Access denied: Insufficient permissions"))
			c.Abort()
			return
		}

		c.Next()
	}
}
