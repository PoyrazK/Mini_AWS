package httphandlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/poyrazk/thecloud/pkg/httputil"
)

// ClusterHandler handles managed Kubernetes HTTP endpoints.
type ClusterHandler struct {
	svc ports.ClusterService
}

// NewClusterHandler constructs a new ClusterHandler.
func NewClusterHandler(svc ports.ClusterService) *ClusterHandler {
	return &ClusterHandler{svc: svc}
}

// CreateClusterRequest is the payload for creating a K8s cluster.
type CreateClusterRequest struct {
	Name    string `json:"name" binding:"required"`
	VpcID   string `json:"vpc_id" binding:"required"`
	Version string `json:"version"`
	Workers int    `json:"workers"`
}

// CreateCluster godoc
// @Summary Create a managed K8s cluster
// @Description Provisions a new Kubernetes cluster using kubeadm
// @Tags K8s
// @Accept json
// @Produce json
// @Param request body CreateClusterRequest true "Cluster details"
// @Success 202 {object} domain.Cluster
// @Router /clusters [post]
func (h *ClusterHandler) CreateCluster(c *gin.Context) {
	var req CreateClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid request body"))
		return
	}

	vpcID, err := uuid.Parse(req.VpcID)
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid vpc_id"))
		return
	}

	userID := appcontext.UserIDFromContext(c.Request.Context())

	// Default version if not specified
	if req.Version == "" {
		req.Version = "v1.29.0"
	}
	if req.Workers == 0 {
		req.Workers = 2 // Default 2 workers
	}

	cluster, err := h.svc.CreateCluster(c.Request.Context(), userID, req.Name, vpcID, req.Version, req.Workers)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusAccepted, cluster)
}

// GetCluster godoc
// @Summary Get cluster details
// @Description Returns cluster metadata and current status
// @Tags K8s
// @Produce json
// @Param id path string true "Cluster ID"
// @Success 200 {object} domain.Cluster
// @Router /clusters/{id} [get]
func (h *ClusterHandler) GetCluster(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cluster id"))
		return
	}

	cluster, err := h.svc.GetCluster(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, cluster)
}

// ListClusters godoc
// @Summary List managed K8s clusters
// @Description Returns all clusters belonging to the user
// @Tags K8s
// @Produce json
// @Success 200 {array} domain.Cluster
// @Router /clusters [get]
func (h *ClusterHandler) ListClusters(c *gin.Context) {
	userID := appcontext.UserIDFromContext(c.Request.Context())
	clusters, err := h.svc.ListClusters(c.Request.Context(), userID)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusOK, clusters)
}

// DeleteCluster godoc
// @Summary Delete a K8s cluster
// @Description Terminates all nodes and removes the cluster record
// @Tags K8s
// @Param id path string true "Cluster ID"
// @Success 202
// @Router /clusters/{id} [delete]
func (h *ClusterHandler) DeleteCluster(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cluster id"))
		return
	}

	if err := h.svc.DeleteCluster(c.Request.Context(), id); err != nil {
		httputil.Error(c, err)
		return
	}

	httputil.Success(c, http.StatusAccepted, nil)
}

// GetKubeconfig godoc
// @Summary Download kubeconfig
// @Description Returns the kubeconfig for clinical access to the cluster
// @Tags K8s
// @Produce plain
// @Param id path string true "Cluster ID"
// @Success 200 {string} string
// @Router /clusters/{id}/kubeconfig [get]
func (h *ClusterHandler) GetKubeconfig(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		httputil.Error(c, errors.New(errors.InvalidInput, "invalid cluster id"))
		return
	}

	kubeconfig, err := h.svc.GetKubeconfig(c.Request.Context(), id)
	if err != nil {
		httputil.Error(c, err)
		return
	}

	c.String(http.StatusOK, kubeconfig)
}
