package services_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFunctionService_Unit(t *testing.T) {
	mockRepo := new(MockFunctionRepository)
	mockCompute := new(MockComputeBackend)
	mockFileStore := new(MockFileStore)
	mockAuditSvc := new(MockAuditService)
	svc := services.NewFunctionService(mockRepo, mockCompute, mockFileStore, mockAuditSvc, slog.Default())
	
	ctx := context.Background()
	userID := uuid.New()
	ctx = appcontext.WithUserID(ctx, userID)

	t.Run("CreateFunction", func(t *testing.T) {
		mockFileStore.On("Write", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(int64(100), nil).Once()
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockAuditSvc.On("Log", mock.Anything, userID, "function.create", "function", mock.Anything, mock.Anything).Return(nil).Once()

		fn, err := svc.CreateFunction(ctx, "test-fn", "nodejs20", "index.handler", []byte("code"))
		assert.NoError(t, err)
		assert.NotNil(t, fn)
		mockRepo.AssertExpectations(t)
	})
}
