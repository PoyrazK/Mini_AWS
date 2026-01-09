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

func TestAutoScalingRepo_ListAllGroups(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewAutoScalingRepo(mock)
		now := time.Now()

		mock.ExpectQuery("(?s)SELECT.*FROM scaling_groups").
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "idempotency_key", "name", "vpc_id", "load_balancer_id", "image", "ports", "min_instances", "max_instances", "desired_count", "current_count", "status", "version", "created_at", "updated_at"}).
				AddRow(uuid.New(), uuid.New(), nil, "group-1", uuid.New(), nil, "image", nil, 1, 10, 2, 2, string(domain.ScalingGroupStatusActive), 1, now, now))

		groups, err := repo.ListAllGroups(context.Background())
		assert.NoError(t, err)
		assert.Len(t, groups, 1)
	})
}

func TestAutoScalingRepo_DeletePolicy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewAutoScalingRepo(mock)
		id := uuid.New()

		mock.ExpectExec("DELETE FROM scaling_policies").
			WithArgs(id).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err = repo.DeletePolicy(context.Background(), id)
		assert.NoError(t, err)
	})
}

func TestAutoScalingRepo_UpdatePolicyLastScaled(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewAutoScalingRepo(mock)
		id := uuid.New()
		now := time.Now()

		mock.ExpectExec("UPDATE scaling_policies").
			WithArgs(pgxmock.AnyArg(), id).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = repo.UpdatePolicyLastScaled(context.Background(), id, now)
		assert.NoError(t, err)
	})
}
