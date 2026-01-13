package services

import (
	"context"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPasswordResetRepository struct {
	mock.Mock
}

func (m *mockPasswordResetRepository) Create(ctx context.Context, token *domain.PasswordResetToken) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *mockPasswordResetRepository) GetByTokenHash(ctx context.Context, hash string) (*domain.PasswordResetToken, error) {
	args := m.Called(ctx, hash)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.PasswordResetToken), args.Error(1)
}

func (m *mockPasswordResetRepository) MarkAsUsed(ctx context.Context, id string) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockPasswordResetRepository) DeleteExpired(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

type mockUserRepository struct {
	mock.Mock
}

func (m *mockUserRepository) Create(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.User), args.Error(1)
}
func (m *mockUserRepository) Update(ctx context.Context, user *domain.User) error {
	return m.Called(ctx, user).Error(0)
}
func (m *mockUserRepository) List(ctx context.Context) ([]*domain.User, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.User), args.Error(1)
}
func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestPasswordResetService_RequestReset(t *testing.T) {
	repo := new(mockPasswordResetRepository)
	userRepo := new(mockUserRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := NewPasswordResetService(repo, userRepo, logger)

	t.Run("success", func(t *testing.T) {
		user := &domain.User{ID: uuid.New(), Email: "test@example.com"}
		userRepo.On("GetByEmail", mock.Anything, "test@example.com").Return(user, nil).Once()
		repo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()

		err := svc.RequestReset(context.Background(), "test@example.com")
		assert.NoError(t, err)
		userRepo.AssertExpectations(t)
		repo.AssertExpectations(t)
	})

	t.Run("user not found returns nil", func(t *testing.T) {
		userRepo.On("GetByEmail", mock.Anything, "none@example.com").Return(nil, assert.AnError).Once()

		err := svc.RequestReset(context.Background(), "none@example.com")
		assert.NoError(t, err)
	})
}

func TestPasswordResetService_ResetPassword(t *testing.T) {
	repo := new(mockPasswordResetRepository)
	userRepo := new(mockUserRepository)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	svc := NewPasswordResetService(repo, userRepo, logger)

	t.Run("success", func(t *testing.T) {
		userID := uuid.New()
		token := &domain.PasswordResetToken{
			ID:        uuid.New(),
			UserID:    userID,
			ExpiresAt: time.Now().Add(time.Hour),
			Used:      false,
		}
		repo.On("MarkAsUsed", mock.Anything, token.ID.String()).Return(nil).Once()
		userRepo.On("GetByID", mock.Anything, userID).Return(&domain.User{ID: userID}, nil).Once()
		userRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		
		// We'd need to know the hash to set expectations precisely, or use mock.Anything
		repo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(token, nil).Once()

		err := svc.ResetPassword(context.Background(), "valid-token", "new-password")
		assert.NoError(t, err)
	})

	t.Run("expired token", func(t *testing.T) {
		token := &domain.PasswordResetToken{
			ExpiresAt: time.Now().Add(-time.Hour),
			Used:      false,
		}
		repo.On("GetByTokenHash", mock.Anything, mock.Anything).Return(token, nil).Once()

		err := svc.ResetPassword(context.Background(), "expired-token", "pass")
		assert.Error(t, err)
		assert.Equal(t, "token expired", err.Error())
	})
}
