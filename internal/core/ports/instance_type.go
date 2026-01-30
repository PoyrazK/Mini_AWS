package ports

import (
	"context"

	"github.com/poyrazk/thecloud/internal/core/domain"
)

// InstanceTypeRepository handles the persistence and retrieval of instance types.
type InstanceTypeRepository interface {
	// List returns all available instance types.
	List(ctx context.Context) ([]*domain.InstanceType, error)
	// GetByID retrieves a specific instance type by its unique identifier.
	GetByID(ctx context.Context, id string) (*domain.InstanceType, error)
}

// InstanceTypeService defines the business logic for instance types.
type InstanceTypeService interface {
	// List returns all available instance types.
	List(ctx context.Context) ([]*domain.InstanceType, error)
}
