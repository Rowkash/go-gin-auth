package middleware

import (
	"errors"
	"io"
	"net/http"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			if appErr, ok := errors.AsType[*common.AppError](err); ok {
				c.AbortWithStatusJSON(appErr.StatusCode, appErr)
				return
			}

			if validationErrors, ok := errors.AsType[validator.ValidationErrors](err); ok {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewValidationError(validationErrors))
				return
			}

			if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
				c.AbortWithStatusJSON(http.StatusBadRequest, common.NewBadRequestError("Invalid request payload structure"))
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, common.NewInternalServerError())
		}
	}
}
