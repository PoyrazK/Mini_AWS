package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poyraz/cloud/internal/core/ports"
	"github.com/poyraz/cloud/pkg/httputil"
)

type IdentityHandler struct {
	svc ports.IdentityService
}

func NewIdentityHandler(svc ports.IdentityService) *IdentityHandler {
	return &IdentityHandler{svc: svc}
}

type CreateKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

func (h *IdentityHandler) CreateKey(c *gin.Context) {
	var req CreateKeyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, err)
		return
	}

	key, err := h.svc.GenerateApiKey(c.Request.Context(), req.Name)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusCreated, key)
}
