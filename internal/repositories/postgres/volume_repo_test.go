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

func TestVolumeRepository_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewVolumeRepository(db)
	ctx := context.Background()

	// Cleanup
	_, err := db.Exec(ctx, "DELETE FROM volumes")
	require.NoError(t, err)

	volID := uuid.New()
	vol := &domain.Volume{
		ID:        volID,
		Name:      "test-vol",
		SizeGB:    10,
		Status:    "available",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	t.Run("Create and Get", func(t *testing.T) {
		err := repo.Create(ctx, vol)
		require.NoError(t, err)

		fetched, err := repo.GetByID(ctx, volID)
		require.NoError(t, err)
		assert.Equal(t, vol.Name, fetched.Name)
	})

	t.Run("GetByName", func(t *testing.T) {
		fetched, err := repo.GetByName(ctx, "test-vol")
		require.NoError(t, err)
		assert.Equal(t, volID, fetched.ID)
	})

	t.Run("List", func(t *testing.T) {
		list, err := repo.List(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, list)
	})

	t.Run("Update and ListByInstanceID", func(t *testing.T) {
		instID := uuid.New()
		// Create instance to satisfy foreign key
		_, err := db.Exec(ctx, "INSERT INTO instances (id, name, image, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			instID, "test-inst", "alpine", "running", 1, time.Now(), time.Now())
		require.NoError(t, err)

		vol.InstanceID = &instID
		vol.MountPath = "/mnt/data"
		vol.UpdatedAt = time.Now()

		err = repo.Update(ctx, vol)
		require.NoError(t, err)

		list, err := repo.ListByInstanceID(ctx, instID)
		require.NoError(t, err)
		assert.Len(t, list, 1)
		assert.Equal(t, volID, list[0].ID)
		assert.Equal(t, "/mnt/data", list[0].MountPath)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, volID)
		require.NoError(t, err)

		_, err = repo.GetByID(ctx, volID)
		assert.Error(t, err)
	})
}
