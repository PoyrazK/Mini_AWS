// Package httphandlers provides HTTP handlers for the API.
package httphandlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/poyrazk/thecloud/pkg/httputil"
)

// CreateRouteRequest define the payload for creating a route.
type CreateRouteRequest struct {
	Name        string   `json:"name" binding:"required"`
	PathPrefix  string   `json:"path_prefix" binding:"required"`
	TargetURL   string   `json:"target_url" binding:"required"`
	Methods     []string `json:"methods"`
	StripPrefix bool     `json:"strip_prefix"`
	RateLimit   int      `json:"rate_limit"`
	Priority    int      `json:"priority"`
}

// GatewayHandler handles API gateway HTTP endpoints.
type GatewayHandler struct {
	svc ports.GatewayService
}

// NewGatewayHandler constructs a GatewayHandler.
func NewGatewayHandler(svc ports.GatewayService) *GatewayHandler {
	return &GatewayHandler{svc: svc}
}

// CreateRoute establishes a new ingress mapping
// @Summary Create a new gateway route
// @Description Registers a new path pattern for the API gateway to proxy to a backend
// @Tags gateway
// @Accept json
// @Produce json
// @Security APIKeyAuth
// @Param request body CreateRouteRequest true "Create route request"
// @Success 201 {object} domain.GatewayRoute
// @Failure 400 {object} httputil.Response
// @Failure 500 {object} httputil.Response
// @Router /gateway/routes [post]
func (h *GatewayHandler) CreateRoute(c *gin.Context) {
	var req CreateRouteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "Invalid request body"))
		return
	}

	if req.RateLimit == 0 {
		req.RateLimit = 100
	}

	params := ports.CreateRouteParams{
		Name:        req.Name,
		Pattern:     req.PathPrefix,
		Target:      req.TargetURL,
		Methods:     req.Methods,
		StripPrefix: req.StripPrefix,
		RateLimit:   req.RateLimit,
		Priority:    req.Priority,
	}

	route, err := h.svc.CreateRoute(c.Request.Context(), params)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusCreated, route)
}

// ListRoutes returns all gateway routes
// @Summary List all gateway routes
// @Description Gets a list of all registered API gateway routes
// @Tags gateway
// @Produce json
// @Security APIKeyAuth
// @Success 200 {array} domain.GatewayRoute
// @Failure 500 {object} httputil.Response
// @Router /gateway/routes [get]
func (h *GatewayHandler) ListRoutes(c *gin.Context) {
	routes, err := h.svc.ListRoutes(c.Request.Context())
	if err != nil {
		httputil.Error(c, err)
		return
	}
	httputil.Success(c, http.StatusOK, routes)
}

// DeleteRoute removes a gateway route
// @Summary Delete a gateway route
// @Description Removes an existing API gateway route by ID
// @Tags gateway
// @Produce json
// @Security APIKeyAuth
// @Param id path string true "Route ID"
// @Success 200 {object} httputil.Response
// @Failure 404 {object} httputil.Response
// @Failure 500 {object} httputil.Response
// @Router /gateway/routes/{id} [delete]
func (h *GatewayHandler) DeleteRoute(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "Invalid route ID"))
		return
	}

	if err := h.svc.DeleteRoute(c.Request.Context(), id); err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, gin.H{"message": "Route deleted"})
}

func (h *GatewayHandler) Proxy(c *gin.Context) {
	path := c.Param("proxy") // Expecting routes like /gw/*proxy
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	proxy, params, ok := h.svc.GetProxy(c.Request.Method, path)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "No route found for " + path})
		return
	}

	// Inject parameters into request context for downstream services if needed
	if len(params) > 0 {
		for k, v := range params {
			c.Set("path_param_"+k, v)
		}
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
