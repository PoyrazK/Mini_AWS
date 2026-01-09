package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestQueueRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresQueueRepository(mock)
	q := &domain.Queue{
		ID:                uuid.New(),
		UserID:            uuid.New(),
		Name:              "test-queue",
		ARN:               "arn:aws:sqs:us-east-1:123456789012:test-queue",
		VisibilityTimeout: 30,
		RetentionDays:     4,
		MaxMessageSize:    262144,
		Status:            domain.QueueStatusActive,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	mock.ExpectExec("INSERT INTO queues").
		WithArgs(q.ID, q.UserID, q.Name, q.ARN, q.VisibilityTimeout, q.RetentionDays, q.MaxMessageSize, q.Status, q.CreatedAt, q.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), q)
	assert.NoError(t, err)
}

func TestQueueRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresQueueRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, arn, visibility_timeout, retention_days, max_message_size, status, created_at, updated_at FROM queues").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "arn", "visibility_timeout", "retention_days", "max_message_size", "status", "created_at", "updated_at"}).
			AddRow(id, userID, "test-queue", "arn", 30, 4, 262144, string(domain.QueueStatusActive), now, now))

	q, err := repo.GetByID(context.Background(), id, userID)
	assert.NoError(t, err)
	assert.NotNil(t, q)
	assert.Equal(t, id, q.ID)
}

func TestQueueRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresQueueRepository(mock)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, arn, visibility_timeout, retention_days, max_message_size, status, created_at, updated_at FROM queues").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "arn", "visibility_timeout", "retention_days", "max_message_size", "status", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "test-queue", "arn", 30, 4, 262144, string(domain.QueueStatusActive), now, now))

	queues, err := repo.List(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, queues, 1)
}
