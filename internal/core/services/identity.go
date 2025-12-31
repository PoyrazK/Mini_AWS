package services

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/poyraz/cloud/internal/core/domain"
	"github.com/poyraz/cloud/internal/core/ports"
	"github.com/poyraz/cloud/internal/errors"
)

type IdentityService struct {
	repo ports.IdentityRepository
}

func NewIdentityService(repo ports.IdentityRepository) *IdentityService {
	return &IdentityService{repo: repo}
}

func (s *IdentityService) GenerateApiKey(ctx context.Context, name string) (*domain.ApiKey, error) {
	// Generate a secure random key
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to generate random key", err)
	}
	keyStr := "miniaws_" + hex.EncodeToString(b)

	apiKey := &domain.ApiKey{
		ID:        uuid.New(),
		UserID:    uuid.New(), // In a real system, this would be the logged-in user
		Key:       keyStr,
		Name:      name,
		CreatedAt: time.Now(),
	}

	if err := s.repo.CreateApiKey(ctx, apiKey); err != nil {
		return nil, err
	}

	return apiKey, nil
}

func (s *IdentityService) ValidateApiKey(ctx context.Context, key string) (bool, error) {
	_, err := s.repo.GetApiKeyByKey(ctx, key)
	if err != nil {
		if errors.Is(err, errors.Unauthorized) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
