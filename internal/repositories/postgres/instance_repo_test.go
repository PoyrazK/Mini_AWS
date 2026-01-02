//go:build integration

package postgres

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB(t *testing.T) *pgxpool.Pool {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://cloud:cloud@localhost:5433/thecloud"
	}

	ctx := context.Background()
	db, err := pgxpool.New(ctx, dbURL)
	require.NoError(t, err)

	err = db.Ping(ctx)
	if err != nil {
		t.Skip("Skipping integration test: database not available")
	}

	return db
}

func TestInstanceRepository_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewInstanceRepository(db)
	ctx := context.Background()

	// Cleanup
	_, err := db.Exec(ctx, "DELETE FROM instances")
	require.NoError(t, err)

	t.Run("Create and Get", func(t *testing.T) {
		id := uuid.New()
		inst := &domain.Instance{
			ID:        id,
			Name:      "integration-test-inst",
			Image:     "alpine",
			Status:    domain.StatusStarting,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			Version:   1,
		}

		err := repo.Create(ctx, inst)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, id)
		require.NoError(t, err)
		assert.Equal(t, inst.Name, fetched.Name)
		assert.Equal(t, inst.Status, fetched.Status)
	})

	t.Run("List", func(t *testing.T) {
		list, err := repo.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, list)
	})

	t.Run("GetByName", func(t *testing.T) {
		fetched, err := repo.GetByName(ctx, "integration-test-inst")
		require.NoError(t, err)
		assert.Equal(t, "integration-test-inst", fetched.Name)
	})

	t.Run("Update", func(t *testing.T) {
		inst, err := repo.GetByName(ctx, "integration-test-inst")
		require.NoError(t, err)

		inst.Status = domain.StatusRunning
		err = repo.Update(ctx, inst)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, inst.ID)
		require.NoError(t, err)
		assert.Equal(t, domain.StatusRunning, fetched.Status)
		assert.Equal(t, 2, fetched.Version)
	})

	t.Run("Update Conflict", func(t *testing.T) {
		inst, err := repo.GetByName(ctx, "integration-test-inst")
		require.NoError(t, err)

		// Create a stale copy
		staleInst := *inst

		// Update original
		inst.Status = domain.StatusStopped
		err = repo.Update(ctx, inst)
		require.NoError(t, err)

		// Try to update with stale copy
		staleInst.Status = domain.StatusStarting
		err = repo.Update(ctx, &staleInst)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflict")
	})

	t.Run("Delete", func(t *testing.T) {
		inst, err := repo.GetByName(ctx, "integration-test-inst")
		require.NoError(t, err)

		err = repo.Delete(ctx, inst.ID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, inst.ID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}
