package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_Register_Unit(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockIdentitySvc := new(MockIdentityService)
	mockAuditSvc := new(MockAuditService)
	mockTenantSvc := new(MockTenantService)

	svc := services.NewAuthService(mockUserRepo, mockIdentitySvc, mockAuditSvc, mockTenantSvc)

	ctx := context.Background()
	email := "test@example.com"
	password := "StrongPass123!@#LongEnoughToPassValidator"
	name := "Test User"

	t.Run("Success", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(nil, nil).Once()
		mockUserRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockTenantSvc.On("CreateTenant", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(&domain.Tenant{ID: uuid.New()}, nil).Once()
		mockUserRepo.On("GetByID", mock.Anything, mock.Anything).Return(&domain.User{ID: uuid.New(), Email: email}, nil).Once()
		mockAuditSvc.On("Log", mock.Anything, mock.Anything, "user.register", "user", mock.Anything, mock.Anything).
			Return(nil).Once()

		user, err := svc.Register(ctx, email, password, name)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, email, user.Email)
		mockUserRepo.AssertExpectations(t)
		mockTenantSvc.AssertExpectations(t)
	})

	t.Run("WeakPassword", func(t *testing.T) {
		user, err := svc.Register(ctx, email, "weak", name)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "weak")
	})

	t.Run("DuplicateEmail", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(&domain.User{ID: uuid.New()}, nil).Once()

		user, err := svc.Register(ctx, email, password, name)

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestAuthService_Login_Unit(t *testing.T) {
	mockUserRepo := new(MockUserRepo)
	mockIdentitySvc := new(MockIdentityService)
	mockAuditSvc := new(MockAuditService)
	mockTenantSvc := new(MockTenantService)

	svc := services.NewAuthService(mockUserRepo, mockIdentitySvc, mockAuditSvc, mockTenantSvc)

	ctx := context.Background()
	email := "login@example.com"
	password := "Password123!"

	t.Run("Success", func(t *testing.T) {
		hp, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		user := &domain.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: string(hp),
		}

		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(user, nil).Once()
		mockIdentitySvc.On("CreateKey", mock.Anything, user.ID, "Default Key").
			Return(&domain.APIKey{Key: "test-key"}, nil).Once()
		mockAuditSvc.On("Log", mock.Anything, user.ID, "user.login", "user", user.ID.String(), mock.Anything).
			Return(nil).Once()

		resultUser, token, err := svc.Login(ctx, email, password)

		assert.NoError(t, err)
		assert.NotNil(t, resultUser)
		assert.Equal(t, "test-key", token)
		mockUserRepo.AssertExpectations(t)
		mockIdentitySvc.AssertExpectations(t)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(nil, nil).Once()

		user, token, err := svc.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
	})

	t.Run("WrongPassword", func(t *testing.T) {
		hp, _ := bcrypt.GenerateFromPassword([]byte("different"), bcrypt.DefaultCost)
		user := &domain.User{
			ID:           uuid.New(),
			Email:        email,
			PasswordHash: string(hp),
		}

		mockUserRepo.On("GetByEmail", mock.Anything, email).Return(user, nil).Once()

		user, token, err := svc.Login(ctx, email, password)

		assert.Error(t, err)
		assert.Empty(t, token)
		assert.Nil(t, user)
	})
}
