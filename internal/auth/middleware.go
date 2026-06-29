package auth

import (
	"strconv"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/Rowkash/go-gin-auth/internal/common/types"
	"github.com/gin-gonic/gin"
)

func middleware(tokenService *TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {

		accessToken, err := c.Cookie("accessToken")

		if err != nil {
			_ = c.Error(common.NewUnauthorizedError("Authorization cookie is required"))
			c.Abort()
			return
		}

		claims, err := tokenService.validate(accessToken)
		if err != nil {
			_ = c.Error(common.NewUnauthorizedError("Invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(types.UserKey, types.UserContext{
			ID:   strconv.Itoa(int(claims.UserId)),
			Role: claims.UserRole,
		})

		c.Next()
	}
}
