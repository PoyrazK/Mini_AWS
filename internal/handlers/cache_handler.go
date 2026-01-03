package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/poyrazk/thecloud/pkg/httputil"
)

type CacheHandler struct {
	svc ports.CacheService
}

func NewCacheHandler(svc ports.CacheService) *CacheHandler {
	return &CacheHandler{svc: svc}
}

type CreateCacheRequest struct {
	Name     string     `json:"name" binding:"required"`
	Version  string     `json:"version" binding:"required"`
	MemoryMB int        `json:"memory_mb" binding:"required"`
	VpcID    *uuid.UUID `json:"vpc_id"`
}

func (h *CacheHandler) Create(c *gin.Context) {
	var req CreateCacheRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, err.Error()))
		return
	}

	cache, err := h.svc.CreateCache(c.Request.Context(), req.Name, req.Version, req.MemoryMB, req.VpcID)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusCreated, cache)
}

func (h *CacheHandler) List(c *gin.Context) {
	caches, err := h.svc.ListCaches(c.Request.Context())
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, caches)
}

func (h *CacheHandler) Get(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cache id"))
		return
	}

	cache, err := h.svc.GetCache(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, cache)
}

func (h *CacheHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cache id"))
		return
	}

	if err := h.svc.DeleteCache(c.Request.Context(), id); err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, gin.H{"message": "cache deleted"})
}

func (h *CacheHandler) GetConnectionString(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cache id"))
		return
	}

	connStr, err := h.svc.GetConnectionString(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, gin.H{"connection_string": connStr})
}

func (h *CacheHandler) Flush(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cache id"))
		return
	}

	if err := h.svc.FlushCache(c.Request.Context(), id); err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, gin.H{"message": "cache flushed"})
}

func (h *CacheHandler) GetStats(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cache id"))
		return
	}

	stats, err := h.svc.GetCacheStats(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, stats)
}
