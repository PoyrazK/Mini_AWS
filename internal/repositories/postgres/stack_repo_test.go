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

func TestStackRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewStackRepository(mock)
	s := &domain.Stack{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Name:       "test-stack",
		Template:   "{}",
		Parameters: []byte(`{"foo": "bar"}`),
		Status:     "CREATE_IN_PROGRESS",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO stacks").
		WithArgs(s.ID, s.UserID, s.Name, s.Template, s.Parameters, s.Status, s.CreatedAt, s.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), s)
	assert.NoError(t, err)
}

func TestStackRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewStackRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, template, parameters, status, status_reason, created_at, updated_at FROM stacks").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "template", "parameters", "status", "status_reason", "created_at", "updated_at"}).
			AddRow(id, userID, "test", "{}", nil, "ACTIVE", nil, now, now))

	mock.ExpectQuery("SELECT id, stack_id, logical_id, physical_id, resource_type, status, created_at FROM stack_resources").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"id", "stack_id", "logical_id", "physical_id", "resource_type", "status", "created_at"}).
			AddRow(uuid.New(), id, "res1", "phys1", "type1", "status1", now))

	s, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Len(t, s.Resources, 1)
}

func TestStackRepository_ListByUserID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewStackRepository(mock)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, template, parameters, status, status_reason, created_at, updated_at FROM stacks").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "template", "parameters", "status", "status_reason", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "s1", "{}", nil, "ACTIVE", nil, now, now))

	stacks, err := repo.ListByUserID(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, stacks, 1)
}

func TestStackRepository_AddResource(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewStackRepository(mock)
	res := &domain.StackResource{
		ID:           uuid.New(),
		StackID:      uuid.New(),
		LogicalID:    "res1",
		PhysicalID:   "phys1",
		ResourceType: "type1",
		Status:       "CREATE_COMPLETE",
		CreatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO stack_resources").
		WithArgs(res.ID, res.StackID, res.LogicalID, res.PhysicalID, res.ResourceType, res.Status, res.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.AddResource(context.Background(), res)
	assert.NoError(t, err)
}
