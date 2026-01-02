//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLBRepository_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewLBRepository(db)
	ctx := context.Background()

	// Cleanup
	_, err := db.Exec(ctx, "DELETE FROM lb_targets")
	require.NoError(t, err)
	_, err = db.Exec(ctx, "DELETE FROM load_balancers")
	require.NoError(t, err)
	_, err = db.Exec(ctx, "DELETE FROM vpcs")
	require.NoError(t, err)
	_, err = db.Exec(ctx, "DELETE FROM instances")
	require.NoError(t, err)

	vpcID := uuid.New()
	// Create VPC
	_, err = db.Exec(ctx, "INSERT INTO vpcs (id, name, network_id, created_at) VALUES ($1, $2, $3, $4)",
		vpcID, "lb-vpc", "net-lb", time.Now())
	require.NoError(t, err)

	lbID := uuid.New()
	lb := &domain.LoadBalancer{
		ID:             lbID,
		Name:           "test-lb",
		VpcID:          vpcID,
		Port:           80,
		Algorithm:      "round-robin",
		Status:         domain.LBStatusCreating,
		IdempotencyKey: "test-key-1",
		Version:        1,
		CreatedAt:      time.Now(),
	}

	t.Run("Create and Get", func(t *testing.T) {
		err := repo.Create(ctx, lb)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, lbID)
		require.NoError(t, err)
		assert.Equal(t, lb.Name, fetched.Name)
		assert.Equal(t, lb.IdempotencyKey, fetched.IdempotencyKey)
	})

	t.Run("GetByIdempotencyKey", func(t *testing.T) {
		fetched, err := repo.GetByIdempotencyKey(ctx, "test-key-1")
		require.NoError(t, err)
		assert.Equal(t, lbID, fetched.ID)
	})

	t.Run("Update", func(t *testing.T) {
		lb.Status = domain.LBStatusActive
		err := repo.Update(ctx, lb)
		require.NoError(t, err)
		assert.Equal(t, 2, lb.Version)

		fetched, err := repo.GetByID(ctx, lbID)
		require.NoError(t, err)
		assert.Equal(t, domain.LBStatusActive, fetched.Status)
	})

	t.Run("Target Management", func(t *testing.T) {
		instID := uuid.New()
		// Create instance first
		_, err := db.Exec(ctx, "INSERT INTO instances (id, name, image, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			instID, "lb-target-inst", "nginx", "running", 1, time.Now(), time.Now())
		require.NoError(t, err)

		target := &domain.LBTarget{
			ID:         uuid.New(),
			LBID:       lbID,
			InstanceID: instID,
			Port:       80,
			Weight:     1,
			Health:     "HEALTHY",
		}

		err = repo.AddTarget(ctx, target)
		require.NoError(t, err)

		targets, err := repo.ListTargets(ctx, lbID)
		require.NoError(t, err)
		assert.Len(t, targets, 1)
		assert.Equal(t, instID, targets[0].InstanceID)

		err = repo.UpdateTargetHealth(ctx, lbID, instID, "UNHEALTHY")
		require.NoError(t, err)

		targets, err = repo.GetTargetsForInstance(ctx, instID)
		require.NoError(t, err)
		assert.Equal(t, "UNHEALTHY", targets[0].Health)

		err = repo.RemoveTarget(ctx, lbID, instID)
		require.NoError(t, err)

		targets, err = repo.ListTargets(ctx, lbID)
		require.NoError(t, err)
		assert.Empty(t, targets)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, lbID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, lbID)
		assert.Error(t, err)
	})
}
