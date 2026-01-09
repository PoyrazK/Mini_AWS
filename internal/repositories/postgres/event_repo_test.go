package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestEventRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewEventRepository(mock)
	e := &domain.Event{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		Action:       "TEST_ACTION",
		ResourceID:   uuid.New().String(),
		ResourceType: "TEST_RES",
		Metadata:     map[string]interface{}{"foo": "bar"},
		CreatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO events").
		WithArgs(e.ID, e.UserID, e.Action, e.ResourceID, e.ResourceType, e.Metadata, e.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), e)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestEventRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewEventRepository(mock)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	limit := 5
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, action, resource_id, resource_type, metadata, created_at FROM events").
		WithArgs(userID, limit).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "action", "resource_id", "resource_type", "metadata", "created_at"}).
			AddRow(uuid.New(), userID, "A1", "RID1", "RT1", nil, now).
			AddRow(uuid.New(), userID, "A2", "RID2", "RT2", nil, now))

	events, err := repo.List(ctx, limit)
	assert.NoError(t, err)
	assert.Len(t, events, 2)
	assert.NoError(t, mock.ExpectationsWereMet())
}
