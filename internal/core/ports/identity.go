package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

type IdentityRepository interface {
	CreateAPIKey(ctx context.Context, apiKey *domain.APIKey) error
	GetAPIKeyByKey(ctx context.Context, key string) (*domain.APIKey, error)
	GetAPIKeyByID(ctx context.Context, id uuid.UUID) (*domain.APIKey, error)
	ListAPIKeysByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.APIKey, error)
	DeleteAPIKey(ctx context.Context, id uuid.UUID) error
}

type IdentityService interface {
	CreateKey(ctx context.Context, userID uuid.UUID, name string) (*domain.APIKey, error)
	ValidateAPIKey(ctx context.Context, key string) (*domain.APIKey, error)
	ListKeys(ctx context.Context, userID uuid.UUID) ([]*domain.APIKey, error)
	RevokeKey(ctx context.Context, userID uuid.UUID, id uuid.UUID) error
	RotateKey(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.APIKey, error)
}
