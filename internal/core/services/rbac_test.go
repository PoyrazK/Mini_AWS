package services_test

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
)

func TestAuthorize(t *testing.T) {
	userRepo := new(MockUserRepo)
	roleRepo := new(MockRoleRepo)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	svc := services.NewRBACService(userRepo, roleRepo, logger)

	ctx := context.Background()
	userID := uuid.New()

	t.Run("Success_ExactPermission", func(t *testing.T) {
		user := &domain.User{ID: userID, Role: "developer"}
		role := &domain.Role{
			Name: "developer",
			Permissions: []domain.Permission{
				domain.PermissionInstanceLaunch,
				domain.PermissionInstanceRead,
			},
		}

		userRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		roleRepo.On("GetRoleByName", ctx, "developer").Return(role, nil).Once()

		err := svc.Authorize(ctx, userID, domain.PermissionInstanceLaunch)
		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
		roleRepo.AssertExpectations(t)
	})

	t.Run("Success_FullAccess", func(t *testing.T) {
		user := &domain.User{ID: userID, Role: "admin"}
		role := &domain.Role{
			Name: "admin",
			Permissions: []domain.Permission{
				domain.PermissionFullAccess,
			},
		}

		userRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		roleRepo.On("GetRoleByName", ctx, "admin").Return(role, nil).Once()

		err := svc.Authorize(ctx, userID, domain.PermissionVpcCreate)
		assert.NoError(t, err)

		userRepo.AssertExpectations(t)
		roleRepo.AssertExpectations(t)
	})

	t.Run("Failure_Denied", func(t *testing.T) {
		user := &domain.User{ID: userID, Role: "viewer"}
		role := &domain.Role{
			Name: "viewer",
			Permissions: []domain.Permission{
				domain.PermissionInstanceRead,
			},
		}

		userRepo.On("GetByID", ctx, userID).Return(user, nil).Once()
		roleRepo.On("GetRoleByName", ctx, "viewer").Return(role, nil).Once()

		err := svc.Authorize(ctx, userID, domain.PermissionInstanceLaunch)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "permission denied")

		userRepo.AssertExpectations(t)
		roleRepo.AssertExpectations(t)
	})
}
