package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

type ContainerRepository interface {
	CreateDeployment(ctx context.Context, d *domain.Deployment) error
	GetDeploymentByID(ctx context.Context, id, userID uuid.UUID) (*domain.Deployment, error)
	ListDeployments(ctx context.Context, userID uuid.UUID) ([]*domain.Deployment, error)
	UpdateDeployment(ctx context.Context, d *domain.Deployment) error
	DeleteDeployment(ctx context.Context, id uuid.UUID) error

	// Replication management
	AddContainer(ctx context.Context, deploymentID, instanceID uuid.UUID) error
	RemoveContainer(ctx context.Context, deploymentID, instanceID uuid.UUID) error
	GetContainers(ctx context.Context, deploymentID uuid.UUID) ([]uuid.UUID, error)

	// Worker
	ListAllDeployments(ctx context.Context) ([]*domain.Deployment, error)
}

type ContainerService interface {
	CreateDeployment(ctx context.Context, name, image string, replicas int, ports string) (*domain.Deployment, error)
	ListDeployments(ctx context.Context) ([]*domain.Deployment, error)
	GetDeployment(ctx context.Context, id uuid.UUID) (*domain.Deployment, error)
	ScaleDeployment(ctx context.Context, id uuid.UUID, replicas int) error
	DeleteDeployment(ctx context.Context, id uuid.UUID) error
}
