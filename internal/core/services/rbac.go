package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

type rbacService struct {
	userRepo ports.UserRepository
	roleRepo ports.RoleRepository
	logger   *slog.Logger
}

func NewRBACService(userRepo ports.UserRepository, roleRepo ports.RoleRepository, logger *slog.Logger) *rbacService {
	return &rbacService{
		userRepo: userRepo,
		roleRepo: roleRepo,
		logger:   logger,
	}
}

func (s *rbacService) Authorize(ctx context.Context, userID uuid.UUID, permission domain.Permission) error {
	allowed, err := s.HasPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !allowed {
		return errors.New(errors.Forbidden, fmt.Sprintf("permission denied: %s", permission))
	}
	return nil
}

func (s *rbacService) HasPermission(ctx context.Context, userID uuid.UUID, permission domain.Permission) (bool, error) {
	// 1. Get user
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to get user: %w", err)
	}

	// 2. Get role
	role, err := s.roleRepo.GetRoleByName(ctx, user.Role)
	if err != nil {
		// If role not found in DB, default to viewer if it was just a string
		s.logger.Warn("role not found in DB, checking if it is a default role", "role", user.Role)
		return false, nil // For now, strict check
	}

	// 3. Check permissions
	for _, p := range role.Permissions {
		if p == domain.PermissionFullAccess {
			return true, nil
		}
		if p == permission {
			return true, nil
		}
	}

	return false, nil
}
