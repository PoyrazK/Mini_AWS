package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/pashagolub/pgxmock/v3"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestAutoScalingRepo_CreateGroup(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewAutoScalingRepo(mock)
	group := &domain.ScalingGroup{
		ID:             uuid.New(),
		UserID:         uuid.New(),
		IdempotencyKey: "key-1",
		Name:           "asg-1",
		VpcID:          uuid.New(),
		LoadBalancerID: nil,
		Image:          "ubuntu",
		Ports:          "80:80",
		MinInstances:   1,
		MaxInstances:   5,
		DesiredCount:   2,
		CurrentCount:   0,
		Status:         domain.ScalingGroupStatusActive,
		Version:        1,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	mock.ExpectExec("INSERT INTO scaling_groups").
		WithArgs(group.ID, group.UserID, group.IdempotencyKey, group.Name, group.VpcID, group.LoadBalancerID,
			group.Image, group.Ports, group.MinInstances, group.MaxInstances,
			group.DesiredCount, group.CurrentCount, group.Status, group.Version,
			group.CreatedAt, group.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateGroup(context.Background(), group)
	assert.NoError(t, err)
}

func TestAutoScalingRepo_GetGroupByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewAutoScalingRepo(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	// Mock SQL.NullString for idk and ports
	// Note: In AddRow, we should use values that match what Scan expects.
	// The repo scans: &g.ID, &g.UserID, &idk, &g.Name, &g.VpcID, &lbID, &g.Image, &ports, ...
	// idk and ports are sql.NullString. lbID is *uuid.UUID.

	idk := sql.NullString{String: "key-1", Valid: true}
	ports := sql.NullString{String: "80:80", Valid: true}
	var lbID *uuid.UUID = nil

	mock.ExpectQuery("SELECT id, user_id, idempotency_key, name, vpc_id, load_balancer_id, image, ports").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "user_id", "idempotency_key", "name", "vpc_id", "load_balancer_id", "image", "ports",
			"min_instances", "max_instances", "desired_count", "current_count", "status", "version", "created_at", "updated_at",
		}).
			AddRow(id, userID, idk, "asg-1", uuid.New(), lbID, "ubuntu", ports,
				1, 5, 2, 0, string(domain.ScalingGroupStatusActive), 1, now, now))

	g, err := repo.GetGroupByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, g)
	assert.Equal(t, id, g.ID)
	assert.Equal(t, "key-1", g.IdempotencyKey)
}

func TestAutoScalingRepo_ListGroups(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewAutoScalingRepo(mock)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	idk := sql.NullString{String: "key-1", Valid: true}
	ports := sql.NullString{String: "80:80", Valid: true}
	var lbID *uuid.UUID = nil

	mock.ExpectQuery("SELECT id, user_id, idempotency_key, name, vpc_id, load_balancer_id, image, ports").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{
			"id", "user_id", "idempotency_key", "name", "vpc_id", "load_balancer_id", "image", "ports",
			"min_instances", "max_instances", "desired_count", "current_count", "status", "version", "created_at", "updated_at",
		}).
			AddRow(uuid.New(), userID, idk, "asg-1", uuid.New(), lbID, "ubuntu", ports,
				1, 5, 2, 0, string(domain.ScalingGroupStatusActive), 1, now, now))

	groups, err := repo.ListGroups(ctx)
	assert.NoError(t, err)
	assert.Len(t, groups, 1)
}
