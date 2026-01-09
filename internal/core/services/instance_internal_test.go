package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

// mockVolumeRepo is already defined in dashboard_test.go (package services)

func TestInstanceService_Internal_GetVolumeByIDOrName(t *testing.T) {
	repo := new(mockVolumeRepo)
	svc := &InstanceService{volumeRepo: repo}
	ctx := context.Background()
	volID := uuid.New()

	t.Run("ByID", func(t *testing.T) {
		repo.On("GetByID", ctx, volID).Return(&domain.Volume{ID: volID}, nil).Once()
		res, err := svc.getVolumeByIDOrName(ctx, volID.String())
		assert.NoError(t, err)
		assert.Equal(t, volID, res.ID)
	})

	t.Run("ByName", func(t *testing.T) {
		repo.On("GetByName", ctx, "test-vol").Return(&domain.Volume{Name: "test-vol"}, nil).Once()
		res, err := svc.getVolumeByIDOrName(ctx, "test-vol")
		assert.NoError(t, err)
		assert.Equal(t, "test-vol", res.Name)
	})
}
