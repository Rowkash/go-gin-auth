package auth

import (
	"net/http"

	"github.com/Rowkash/go-gin-auth/internal/common"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) register(c *gin.Context) {
	var body RegisterDto

	if err := c.ShouldBind(&body); err != nil {
		_ = c.Error(err)
		return
	}
	tokens, err := h.service.register(c.Request.Context(), &body)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SetCookie(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, "Successfully registered")
}

func (h *Handler) login(c *gin.Context) {
	var body LoginDto

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokens, err := h.service.login(c.Request.Context(), &body)
	if err != nil {
		_ = c.Error(err)
		return
	}
	SetCookie(c, tokens.AccessToken, tokens.RefreshToken)

	c.JSON(http.StatusOK, "Successfully logged in")
}

func (h *Handler) logout(c *gin.Context) {
	refreshToken, err := c.Cookie("refreshToken")
	if err != nil {
		_ = c.Error(common.NewBadRequestError("Refresh token cookie is required to logout"))
		c.Abort()
		return
	}

	err = h.service.logout(c.Request.Context(), refreshToken)
	if err != nil {
		_ = c.Error(err)
		return
	}
	ClearCookie(c)
	c.JSON(http.StatusOK, "Successfully logged out")
}
