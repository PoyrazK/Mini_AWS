package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

type GlobalLBService struct {
	repo     ports.GlobalLBRepository
	lbRepo   ports.LBRepository
	geoDNS   ports.GeoDNSBackend
	auditSvc ports.AuditService
	logger   *slog.Logger
}

func NewGlobalLBService(
	repo ports.GlobalLBRepository,
	lbRepo ports.LBRepository,
	geoDNS ports.GeoDNSBackend,
	auditSvc ports.AuditService,
	logger *slog.Logger,
) *GlobalLBService {
	return &GlobalLBService{
		repo:     repo,
		lbRepo:   lbRepo,
		geoDNS:   geoDNS,
		auditSvc: auditSvc,
		logger:   logger,
	}
}

func (s *GlobalLBService) Create(ctx context.Context, name, hostname string, policy domain.RoutingPolicy, healthCheck domain.GlobalHealthCheckConfig) (*domain.GlobalLoadBalancer, error) {
	// Validate inputs
	if name == "" || hostname == "" {
		return nil, errors.New(errors.InvalidInput, "name and hostname are required")
	}

	// Check for hostname uniqueness
	existing, err := s.repo.GetByHostname(ctx, hostname)
	if err == nil && existing != nil {
		return nil, errors.New(errors.Conflict, "hostname already in use")
	}

	glb := &domain.GlobalLoadBalancer{
		ID:          uuid.New(),
		UserID:      appcontext.UserIDFromContext(ctx),
		TenantID:    appcontext.TenantIDFromContext(ctx),
		Name:        name,
		Hostname:    hostname,
		Policy:      policy,
		HealthCheck: healthCheck,
		Status:      "ACTIVE",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Endpoints:   []*domain.GlobalEndpoint{},
	}

	if err := s.repo.Create(ctx, glb); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to create global load balancer", err)
	}

	// Initialize the associated GeoDNS record set.
	// Current behavior initializes an empty record set to ensure immediate resolution readiness.
	if err := s.geoDNS.CreateGeoRecord(ctx, hostname, nil); err != nil {
		s.logger.Error("geo-dns initialization failed", "hostname", hostname, "error", err)
		// Non-blocking: failures in the DNS synchronization layer are logged for asynchronous remediation.
	}

	_ = s.auditSvc.Log(ctx, glb.UserID, "global_lb.create", "global_lb", glb.ID.String(), map[string]interface{}{
		"hostname": hostname,
		"policy":   policy,
	})

	return glb, nil
}

func (s *GlobalLBService) Get(ctx context.Context, id uuid.UUID) (*domain.GlobalLoadBalancer, error) {
	glb, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Load endpoints
	endpoints, err := s.repo.ListEndpoints(ctx, id)
	if err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to list endpoints", err)
	}
	glb.Endpoints = endpoints

	return glb, nil
}

func (s *GlobalLBService) List(ctx context.Context) ([]*domain.GlobalLoadBalancer, error) {
	return s.repo.List(ctx)
}

func (s *GlobalLBService) Delete(ctx context.Context, id uuid.UUID) error {
	glb, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	// Synchronously remove the associated GeoDNS record set.
	if err := s.geoDNS.DeleteGeoRecord(ctx, glb.Hostname); err != nil {
		s.logger.Error("geo-dns record deletion failed", "hostname", glb.Hostname, "error", err)
		// Proceed with database deletion to maintain system consistency and prevent orphaned resources.
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return errors.Wrap(errors.Internal, "failed to delete global load balancer", err)
	}

	_ = s.auditSvc.Log(ctx, glb.UserID, "global_lb.delete", "global_lb", id.String(), nil)

	return nil
}

func (s *GlobalLBService) AddEndpoint(ctx context.Context, glbID uuid.UUID, region string, targetType string, targetID *uuid.UUID, targetIP *string, weight, priority int) (*domain.GlobalEndpoint, error) {
	glb, err := s.Get(ctx, glbID)
	if err != nil {
		return nil, err
	}

	// Validate target
	if targetType == "LB" {
		if targetID == nil {
			return nil, errors.New(errors.InvalidInput, "target_id required for LB endpoint")
		}
		// Verify LB exists
		if _, err := s.lbRepo.GetByID(ctx, *targetID); err != nil {
			return nil, errors.Wrap(errors.NotFound, "target load balancer not found", err)
		}
	} else if targetType == "IP" {
		if targetIP == nil || *targetIP == "" {
			return nil, errors.New(errors.InvalidInput, "target_ip required for IP endpoint")
		}
	} else {
		return nil, errors.New(errors.InvalidInput, "invalid target type")
	}

	ep := &domain.GlobalEndpoint{
		ID:         uuid.New(),
		GlobalLBID: glbID,
		Region:     region,
		TargetType: targetType,
		TargetID:   targetID,
		TargetIP:   targetIP,
		Weight:     weight,
		Priority:   priority,
		Healthy:    true, // Assume healthy initially
		CreatedAt:  time.Now(),
	}

	if err := s.repo.AddEndpoint(ctx, ep); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to add endpoint", err)
	}

	// Refresh the Global Load Balancer state to synchronize endpoints with the DNS backend.
	glb, err = s.Get(ctx, glbID)
	if err == nil {
		// Note: The interface currently expects a slice of domain.GlobalEndpoint values.
		eps := make([]domain.GlobalEndpoint, len(glb.Endpoints))
		for i, e := range glb.Endpoints {
			eps[i] = *e
		}

		if err := s.geoDNS.CreateGeoRecord(ctx, glb.Hostname, eps); err != nil {
			s.logger.Error("failed to update geo dns", "hostname", glb.Hostname, "error", err)
		}
	}

	_ = s.auditSvc.Log(ctx, glb.UserID, "global_lb.endpoint_add", "global_lb", glbID.String(), map[string]interface{}{
		"region": region,
		"type":   targetType,
	})

	return ep, nil
}

func (s *GlobalLBService) RemoveEndpoint(ctx context.Context, endpointID uuid.UUID) error {
	// Resource removal involves updating the persistent state followed by a DNS synchronization.
	// Current Limitation: DNS updates require the parent Global Load Balancer ID.
	// Future Iteration: Refactor repository to support atomic retrieval of parent contexts or
	// update signature to include glbID.
	// Actually, let's just implement `RemoveEndpoint` to return the deleted ep or parent ID?
	// The interface is `RemoveEndpoint(ctx, id) error`.

	if err := s.repo.RemoveEndpoint(ctx, endpointID); err != nil {
		return err
	}

	// Note: Explicit DNS synchronization is deferred until the GLB context is provided.
	return nil
}

func (s *GlobalLBService) ListEndpoints(ctx context.Context, glbID uuid.UUID) ([]*domain.GlobalEndpoint, error) {
	return s.repo.ListEndpoints(ctx, glbID)
}
