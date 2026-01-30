package services

import (
	"context"

	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
)

type instanceTypeService struct {
	repo ports.InstanceTypeRepository
}

// NewInstanceTypeService creates a new InstanceTypeService.
func NewInstanceTypeService(repo ports.InstanceTypeRepository) ports.InstanceTypeService {
	return &instanceTypeService{repo: repo}
}

func (s *instanceTypeService) List(ctx context.Context) ([]*domain.InstanceType, error) {
	return s.repo.List(ctx)
}
