package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
)

type ContainerService struct {
	repo     ports.ContainerRepository
	eventSvc ports.EventService
}

func NewContainerService(repo ports.ContainerRepository, eventSvc ports.EventService) ports.ContainerService {
	return &ContainerService{
		repo:     repo,
		eventSvc: eventSvc,
	}
}

func (s *ContainerService) CreateDeployment(ctx context.Context, name, image string, replicas int, ports string) (*domain.Deployment, error) {
	userID := appcontext.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}

	dep := &domain.Deployment{
		ID:           uuid.New(),
		UserID:       userID,
		Name:         name,
		Image:        image,
		Replicas:     replicas,
		CurrentCount: 0,
		Ports:        ports,
		Status:       domain.DeploymentStatusScaling,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.repo.CreateDeployment(ctx, dep); err != nil {
		return nil, err
	}

	_ = s.eventSvc.RecordEvent(ctx, "DEPLOYMENT_CREATED", dep.ID.String(), "DEPLOYMENT", nil)

	return dep, nil
}

func (s *ContainerService) ListDeployments(ctx context.Context) ([]*domain.Deployment, error) {
	userID := appcontext.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.repo.ListDeployments(ctx, userID)
}

func (s *ContainerService) GetDeployment(ctx context.Context, id uuid.UUID) (*domain.Deployment, error) {
	userID := appcontext.UserIDFromContext(ctx)
	if userID == uuid.Nil {
		return nil, fmt.Errorf("unauthorized")
	}
	return s.repo.GetDeploymentByID(ctx, id, userID)
}

func (s *ContainerService) ScaleDeployment(ctx context.Context, id uuid.UUID, replicas int) error {
	userID := appcontext.UserIDFromContext(ctx)
	dep, err := s.repo.GetDeploymentByID(ctx, id, userID)
	if err != nil {
		return err
	}

	dep.Replicas = replicas
	dep.Status = domain.DeploymentStatusScaling
	return s.repo.UpdateDeployment(ctx, dep)
}

func (s *ContainerService) DeleteDeployment(ctx context.Context, id uuid.UUID) error {
	userID := appcontext.UserIDFromContext(ctx)
	dep, err := s.repo.GetDeploymentByID(ctx, id, userID)
	if err != nil {
		return err
	}

	dep.Status = domain.DeploymentStatusDeleting
	return s.repo.UpdateDeployment(ctx, dep)
}
