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

func TestLBRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewLBRepository(mock)
	lb := &domain.LoadBalancer{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		IdempotencyKey: "key-1",
		Name:           "lb-1",
		VpcID:          uuid.New(),
		Port:           80,
		Algorithm:      "round-robin",
		Status:         domain.LBStatusActive,
		Version:        1,
		CreatedAt:      time.Now(),
	}

	mock.ExpectExec("INSERT INTO load_balancers").
		WithArgs(lb.ID, lb.UserID, lb.IdempotencyKey, lb.Name, lb.VpcID, lb.Port, lb.Algorithm, lb.Status, lb.Version, lb.CreatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), lb)
	assert.NoError(t, err)
}

func TestLBRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewLBRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, COALESCE\\(idempotency_key, ''\\), name, vpc_id, port, algorithm, status, version, created_at FROM load_balancers").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "idempotency_key", "name", "vpc_id", "port", "algorithm", "status", "version", "created_at"}).
			AddRow(id, userID, "key-1", "lb-1", uuid.New(), 80, "round-robin", string(domain.LBStatusActive), 1, now))

	lb, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	if lb != nil {
		assert.Equal(t, id, lb.ID)
	}
}

func TestLBRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewLBRepository(mock)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	// List query handles COALESCE(idempotency_key, '')
	mock.ExpectQuery("SELECT id, user_id, COALESCE\\(idempotency_key, ''\\), name, vpc_id, port, algorithm, status, version, created_at FROM load_balancers").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "idempotency_key", "name", "vpc_id", "port", "algorithm", "status", "version", "created_at"}).
			AddRow(uuid.New(), userID, "key-1", "lb-1", uuid.New(), 80, "round-robin", string(domain.LBStatusActive), 1, now))

	lbs, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, lbs, 1)
}
