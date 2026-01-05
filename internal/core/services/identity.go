package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

type IdentityService struct {
	repo ports.IdentityRepository
}

func NewIdentityService(repo ports.IdentityRepository) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) CreateKey(ctx context.Context, userID uuid.UUID, name string) (*domain.APIKey, error) {
	// Generate a secure random key
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to generate random key", err)
	}
	keyStr := "thecloud_" + hex.EncodeToString(b)

	apiKey := &domain.APIKey{
		ID:        uuid.New(),
		UserID:    userID,
		Key:       keyStr,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateAPIKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *IdentityService) ValidateAPIKey(ctx context.Context, key string) (*domain.APIKey, error) {
	apiKey, err := s.repo.GetAPIKeyByKey(ctx, key)
	if err != nil {
		return nil, err
	}
	return apiKey, nil
}

func (s *IdentityService) ListKeys(ctx context.Context, userID uuid.UUID) ([]*domain.APIKey, error) {
	return s.repo.ListAPIKeysByUserID(ctx, userID)
}

func (s *IdentityService) RevokeKey(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	key, err := s.repo.GetAPIKeyByID(ctx, id)
	if err != nil {
		return err
	}

	if key.UserID != userID {
		return errors.New(errors.Forbidden, "cannot revoke key owned by another user")
	}

	return s.repo.DeleteAPIKey(ctx, id)
}

func (s *IdentityService) RotateKey(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.APIKey, error) {
	key, err := s.repo.GetAPIKeyByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if key.UserID != userID {
		return nil, errors.New(errors.Forbidden, "cannot rotate key owned by another user")
	}

	// Create new key
	newKey, err := s.CreateKey(ctx, userID, key.Name+" (rotated)")
	if err != nil {
		return nil, err
	}

	// Delete old key
	if err := s.repo.DeleteAPIKey(ctx, id); err != nil {
		// Log error but we already have a new key
		return newKey, nil
	}

	return newKey, nil
}
