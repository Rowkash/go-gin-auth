package users

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(s *Service) *Handler {
	return &Handler{service: s}
}

func (h *Handler) findById(c *gin.Context) {
	var param FindByIdParams

	if err := c.ShouldBindUri(&param); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.FindOne(UserFilter{ID: &param.Id})
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *Handler) getPage(c *gin.Context) {
	var query AdminUsersPageRequest

	if err := c.ShouldBindQuery(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	users, err := h.service.getPage(query)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.JSON(http.StatusOK, users)
}
