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

func TestInstanceRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewInstanceRepository(mock)
	inst := &domain.Instance{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "inst-1",
		Image:       "ubuntu:latest",
		ContainerID: "cid-1",
		Status:      domain.StatusRunning,
		Ports:       "80:80",
		VpcID:       nil,
		SubnetID:    nil,
		PrivateIP:   "10.0.0.1",
		OvsPort:     "ovs-1",
		Version:     1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO instances").
		WithArgs(inst.ID, inst.UserID, inst.Name, inst.Image, inst.ContainerID, inst.Status, inst.Ports, inst.VpcID, inst.SubnetID, inst.PrivateIP, inst.OvsPort, inst.Version, inst.CreatedAt, inst.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), inst)
	assert.NoError(t, err)
}

func TestInstanceRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewInstanceRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, image, COALESCE\\(container_id, ''\\), status, COALESCE\\(ports, ''\\), vpc_id, subnet_id, COALESCE\\(private_ip::text, ''\\), COALESCE\\(ovs_port, ''\\), version, created_at, updated_at").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "container_id", "status", "ports", "vpc_id", "subnet_id", "private_ip", "ovs_port", "version", "created_at", "updated_at"}).
			AddRow(id, userID, "inst-1", "ubuntu:latest", "cid-1", string(domain.StatusRunning), "80:80", nil, nil, "10.0.0.1", "ovs-1", 1, now, now))

	inst, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, inst)
	assert.Equal(t, id, inst.ID)
}

func TestInstanceRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewInstanceRepository(mock)
	inst := &domain.Instance{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "inst-1-updated",
		Status:      domain.StatusStopped,
		ContainerID: "cid-1",
		Ports:       "80:80",
		Version:     1,
		UpdatedAt:   time.Now(),
	}

	// Update expects optimistic locking check (rows affected)
	mock.ExpectExec("UPDATE instances").
		WithArgs(inst.Name, inst.Status, pgxmock.AnyArg(), inst.ContainerID, inst.Ports, inst.VpcID, inst.SubnetID, inst.PrivateIP, inst.OvsPort, inst.ID, inst.Version, inst.UserID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), inst)
	assert.NoError(t, err)
	assert.Equal(t, 2, inst.Version) // Version should be incremented
}
