package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/mock"
)

type TaskQueueStub struct{}

func (q *TaskQueueStub) Enqueue(ctx context.Context, queueName string, payload interface{}) error {
	return nil
}

func (q *TaskQueueStub) Dequeue(ctx context.Context, queueName string) (string, error) {
	return "", nil
}

type mockEventService struct {
	mock.Mock
}

func (m *mockEventService) RecordEvent(ctx context.Context, action, resourceID, resourceType string, metadata map[string]interface{}) error {
	args := m.Called(ctx, action, resourceID, resourceType, metadata)
	return args.Error(0)
}

func (m *mockEventService) ListEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]*domain.Event), args.Error(1)
}

type mockAuditService struct {
	mock.Mock
}

func (m *mockAuditService) Log(ctx context.Context, userID uuid.UUID, action, resourceType, resourceID string, details map[string]interface{}) error {
	args := m.Called(ctx, userID, action, resourceType, resourceID, details)
	return args.Error(0)
}

func (m *mockAuditService) ListLogs(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]*domain.AuditLog), args.Error(1)
}
