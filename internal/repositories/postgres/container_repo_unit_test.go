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

func TestContainerRepository_CreateDeployment(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresContainerRepository(mock)
	deployment := &domain.Deployment{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		Name:         "test-dep",
		Image:        "nginx",
		Replicas:     3,
		CurrentCount: 0,
		Ports:        "80:80",
		Status:       domain.DeploymentStatusScaling,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	mock.ExpectExec("INSERT INTO deployments").
		WithArgs(deployment.ID, deployment.UserID, deployment.Name, deployment.Image, deployment.Replicas, deployment.CurrentCount, deployment.Ports, deployment.Status, deployment.CreatedAt, deployment.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateDeployment(context.Background(), deployment)
	assert.NoError(t, err)
}

func TestContainerRepository_GetDeploymentByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresContainerRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, image, replicas, current_count, ports, status, created_at, updated_at FROM deployments").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "replicas", "current_count", "ports", "status", "created_at", "updated_at"}).
			AddRow(id, userID, "test-dep", "nginx", 3, 0, "80:80", string(domain.DeploymentStatusScaling), now, now))

	d, err := repo.GetDeploymentByID(context.Background(), id, userID)
	assert.NoError(t, err)
	assert.NotNil(t, d)
	assert.Equal(t, id, d.ID)
	assert.Equal(t, domain.DeploymentStatusScaling, d.Status)
}

func TestContainerRepository_ListDeployments(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresContainerRepository(mock)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, image, replicas, current_count, ports, status, created_at, updated_at FROM deployments").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "replicas", "current_count", "ports", "status", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "test-dep", "nginx", 3, 0, "80:80", string(domain.DeploymentStatusScaling), now, now))

	deps, err := repo.ListDeployments(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, deps, 1)
	assert.Equal(t, domain.DeploymentStatusScaling, deps[0].Status)
}
