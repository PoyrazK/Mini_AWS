package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

// ClusterService implements the managed Kubernetes service.
type ClusterService struct {
	repo        ports.ClusterRepository
	provisioner ports.ClusterProvisioner
	vpcSvc      ports.VpcService
	instanceSvc ports.InstanceService
	logger      *slog.Logger
}

// ClusterServiceParams holds dependencies for ClusterService.
type ClusterServiceParams struct {
	Repo        ports.ClusterRepository
	Provisioner ports.ClusterProvisioner
	VpcSvc      ports.VpcService
	InstanceSvc ports.InstanceService
	Logger      *slog.Logger
}

// NewClusterService constructs a new ClusterService.
func NewClusterService(params ClusterServiceParams) *ClusterService {
	return &ClusterService{
		repo:        params.Repo,
		provisioner: params.Provisioner,
		vpcSvc:      params.VpcSvc,
		instanceSvc: params.InstanceSvc,
		logger:      params.Logger,
	}
}

// CreateCluster initiates the provisioning of a new Kubernetes cluster.
func (s *ClusterService) CreateCluster(ctx context.Context, userID uuid.UUID, name string, vpcID uuid.UUID, version string, workers int) (*domain.Cluster, error) {
	// 1. Verify VPC exists and belongs to user
	vpc, err := s.vpcSvc.GetVPC(ctx, vpcID.String())
	if err != nil {
		return nil, errors.Wrap(errors.NotFound, "vpc not found", err)
	}

	// 2. Create cluster record in database
	cluster := &domain.Cluster{
		ID:          uuid.New(),
		Name:        name,
		UserID:      userID,
		VpcID:       vpc.ID,
		Version:     version,
		WorkerCount: workers,
		Status:      domain.ClusterStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.repo.Create(ctx, cluster); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to create cluster record", err)
	}

	// 3. Start provisioning workflow asynchronously
	go func() {
		// Use a detached context for long-running provisioning
		bgCtx := context.Background()

		s.logger.Info("starting cluster provisioning", "cluster_id", cluster.ID, "name", cluster.Name)

		cluster.Status = domain.ClusterStatusProvisioning
		cluster.UpdatedAt = time.Now()
		_ = s.repo.Update(bgCtx, cluster)

		if err := s.provisioner.Provision(bgCtx, cluster); err != nil {
			s.logger.Error("cluster provisioning failed", "cluster_id", cluster.ID, "error", err)
			cluster.Status = domain.ClusterStatusFailed
		} else {
			s.logger.Info("cluster provisioning completed", "cluster_id", cluster.ID)
			cluster.Status = domain.ClusterStatusRunning
		}

		cluster.UpdatedAt = time.Now()
		_ = s.repo.Update(bgCtx, cluster)
	}()

	return cluster, nil
}

// GetCluster retrieves cluster details by ID.
func (s *ClusterService) GetCluster(ctx context.Context, id uuid.UUID) (*domain.Cluster, error) {
	cluster, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, errors.New(errors.NotFound, "cluster not found")
	}
	return cluster, nil
}

// ListClusters retrieves all clusters for a user.
func (s *ClusterService) ListClusters(ctx context.Context, userID uuid.UUID) ([]*domain.Cluster, error) {
	return s.repo.ListByUserID(ctx, userID)
}

// DeleteCluster removes a cluster and its associated resources.
func (s *ClusterService) DeleteCluster(ctx context.Context, id uuid.UUID) error {
	cluster, err := s.GetCluster(ctx, id)
	if err != nil {
		return err
	}

	cluster.Status = domain.ClusterStatusDeleting
	cluster.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, cluster); err != nil {
		return errors.Wrap(errors.Internal, "failed to update cluster status", err)
	}

	go func() {
		bgCtx := context.Background()
		s.logger.Info("starting cluster deprovisioning", "cluster_id", cluster.ID)

		if err := s.provisioner.Deprovision(bgCtx, cluster); err != nil {
			s.logger.Error("cluster deprovisioning failed", "cluster_id", cluster.ID, "error", err)
			cluster.Status = domain.ClusterStatusFailed
			cluster.UpdatedAt = time.Now()
			_ = s.repo.Update(bgCtx, cluster)
			return
		}

		if err := s.repo.Delete(bgCtx, cluster.ID); err != nil {
			s.logger.Error("failed to delete cluster record", "cluster_id", cluster.ID, "error", err)
		}
		s.logger.Info("cluster deprovisioning completed", "cluster_id", cluster.ID)
	}()

	return nil
}

// GetKubeconfig retrieves the encrypted kubeconfig for a cluster.
func (s *ClusterService) GetKubeconfig(ctx context.Context, id uuid.UUID) (string, error) {
	cluster, err := s.GetCluster(ctx, id)
	if err != nil {
		return "", err
	}

	if cluster.Status != domain.ClusterStatusRunning {
		return "", errors.New(errors.Validation, "kubeconfig is only available when cluster is running")
	}

	if cluster.Kubeconfig == "" {
		return "", errors.New(errors.NotFound, "kubeconfig not found for cluster")
	}

	return cluster.Kubeconfig, nil
}
