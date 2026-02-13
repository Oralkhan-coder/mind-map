package transport

import (
	"net/http"

	"github.com/Oralkhan-coder/mind-map/internal/core"
	"github.com/Oralkhan-coder/mind-map/internal/http/dto"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	srv AuthService
}

func NewAuthHandler(srv AuthService) *AuthHandler {
	return &AuthHandler{
		srv: srv,
	}
}

func (h *AuthHandler) SignUp(c *gin.Context) {
	var req dto.SignUpRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	id, err := h.srv.SignUp(c, req)
	if err != nil {
		c.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *AuthHandler) ConfirmEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		c.Error(core.BadRequest("token required"))
		return
	}

	if err := h.srv.ConfirmEmail(c.Request.Context(), token); err != nil {
		c.Error(err)
		return
	}

	c.Status(http.StatusOK)
}
