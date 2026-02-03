//go:build integration

package services

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/poyrazk/thecloud/internal/repositories/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// StubGeoDNSBackend implements ports.GeoDNSBackend for testing purposes.
type StubGeoDNSBackend struct {
	CreatedRecords map[string][]domain.GlobalEndpoint
	DeletedRecords []string
}

func (s *StubGeoDNSBackend) CreateGeoRecord(ctx context.Context, hostname string, endpoints []domain.GlobalEndpoint) error {
	if s.CreatedRecords == nil {
		s.CreatedRecords = make(map[string][]domain.GlobalEndpoint)
	}
	s.CreatedRecords[hostname] = endpoints
	return nil
}

func (s *StubGeoDNSBackend) DeleteGeoRecord(ctx context.Context, hostname string) error {
	s.DeletedRecords = append(s.DeletedRecords, hostname)
	return nil
}

// MockAuditService implements ports.AuditService for testing purposes.
type MockAuditService struct{}

func (m *MockAuditService) Log(ctx context.Context, userID uuid.UUID, action, resourceType, resourceID string, details map[string]interface{}) error {
	return nil
}

func (m *MockAuditService) ListLogs(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	return nil, nil
}

func TestGlobalLBServiceIntegration(t *testing.T) {
	db := postgres.SetupDB(t)
	defer db.Close()

	repo := postgres.NewGlobalLBRepository(db)
	lbRepo := postgres.NewLBRepository(db) // Corrected constructor name
	geoDNS := &StubGeoDNSBackend{}
	auditSvc := &MockAuditService{}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	svc := NewGlobalLBService(
		repo,
		lbRepo,
		geoDNS,
		auditSvc,
		logger,
	)

	t.Run("Scenario 1: Enforce Hostname Uniqueness", func(t *testing.T) {
		postgres.CleanDB(t, db)
		ctx := postgres.SetupTestUser(t, db)
		hostname := "unique.global.com"

		// 1. Create first GLB
		_, err := svc.Create(ctx, "glb-1", hostname, domain.RoutingLatency, domain.GlobalHealthCheckConfig{Protocol: "HTTP", Port: 80})
		require.NoError(t, err)

		// 2. Try creating second with same hostname
		_, err = svc.Create(ctx, "glb-2", hostname, domain.RoutingWeighted, domain.GlobalHealthCheckConfig{Protocol: "HTTP", Port: 80})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.Conflict))
	})

	t.Run("Scenario 2: Endpoint Default Values", func(t *testing.T) {
		postgres.CleanDB(t, db)
		ctx := postgres.SetupTestUser(t, db)

		glb, err := svc.Create(ctx, "defaults-test", "defaults.global.com", domain.RoutingLatency, domain.GlobalHealthCheckConfig{Protocol: "HTTP", Port: 80})
		require.NoError(t, err)

		// Add endpoint with 0 weight/priority (which should be defaulted by service if logic is there)
		// Note: The current service implementation takes weight/priority as params.
		// If they are 0, we can check if domain/DB layer handles them or if service should.
		ip := "1.1.1.1"
		ep, err := svc.AddEndpoint(ctx, glb.ID, "us-east-1", "IP", nil, &ip, 0, 0)
		require.NoError(t, err)

		// Verify what was actually saved
		endpoints, err := svc.ListEndpoints(ctx, glb.ID)
		require.NoError(t, err)
		assert.Len(t, endpoints, 1)

		// If our business rule says 0 is allowed, they stay 0.
		// If we want defaults, we check for them here.
		assert.Equal(t, 0, ep.Weight)
	})

	t.Run("Scenario 3: State Consistency - Delete Cascade", func(t *testing.T) {
		postgres.CleanDB(t, db)
		ctx := postgres.SetupTestUser(t, db)

		glb, err := svc.Create(ctx, "cascade-test", "cascade.global.com", domain.RoutingLatency, domain.GlobalHealthCheckConfig{Protocol: "HTTP", Port: 80})
		require.NoError(t, err)

		ip := "2.2.2.2"
		_, err = svc.AddEndpoint(ctx, glb.ID, "us-west-1", "IP", nil, &ip, 100, 1)
		require.NoError(t, err)

		// Delete GLB
		err = svc.Delete(ctx, glb.ID)
		require.NoError(t, err)

		// Verify GLB is gone
		_, err = svc.Get(ctx, glb.ID)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.NotFound))

		// Verify endpoints are gone (Cascading delete in DB)
		eps, err := svc.ListEndpoints(ctx, glb.ID)
		assert.NoError(t, err) // Service ListEndpoints might return empty slice
		assert.Empty(t, eps)
	})

	t.Run("Scenario 4: Resilience - Database Disconnection", func(t *testing.T) {
		postgres.CleanDB(t, db)
		ctx := postgres.SetupTestUser(t, db)

		// Close DB to simulate failure
		db.Close()

		_, err := svc.Create(ctx, "fail-test", "fail.global.com", domain.RoutingLatency, domain.GlobalHealthCheckConfig{Protocol: "HTTP", Port: 80})

		assert.Error(t, err)
		assert.True(t, errors.Is(err, errors.Internal))

		// Need to reopen DB for other tests if this wasn't the last one
		// But in this suite it's okay.
	})
}
