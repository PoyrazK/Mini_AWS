package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/pkg/httputil"
)

// InstanceTypeHandler handles instance type related HTTP requests.
type InstanceTypeHandler struct {
	svc ports.InstanceTypeService
}

// NewInstanceTypeHandler creates a new InstanceTypeHandler.
func NewInstanceTypeHandler(svc ports.InstanceTypeService) *InstanceTypeHandler {
	return &InstanceTypeHandler{svc: svc}
}

// List returns all available instance types.
// @Summary List instance types
// @Description Gets a list of all available instance types with their resource limits and pricing
// @Tags instances
// @Produce json
// @Security APIKeyAuth
// @Success 200 {array} domain.InstanceType
// @Failure 500 {object} httputil.Response
// @Router /instance-types [get]
func (h *InstanceTypeHandler) List(c *gin.Context) {
	types, err := h.svc.List(c.Request.Context())
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, types)
}
