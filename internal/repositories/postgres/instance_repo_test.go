package postgres

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	theclouderrors "github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
)

func TestInstanceRepository_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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

		mock.ExpectExec("(?s)INSERT INTO instances.*").
			WithArgs(inst.ID, inst.UserID, inst.Name, inst.Image, inst.ContainerID, string(inst.Status), inst.Ports, inst.VpcID, inst.SubnetID, inst.PrivateIP, inst.OvsPort, inst.Version, inst.CreatedAt, inst.UpdatedAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err = repo.Create(context.Background(), inst)
		assert.NoError(t, err)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		inst := &domain.Instance{ID: uuid.New()}

		mock.ExpectExec("(?s)INSERT INTO instances.*").
			WillReturnError(errors.New("db error"))

		err = repo.Create(context.Background(), inst)
		assert.Error(t, err)
	})
}

func TestInstanceRepository_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(id, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "container_id", "status", "ports", "vpc_id", "subnet_id", "private_ip", "ovs_port", "version", "created_at", "updated_at"}).
				AddRow(id, userID, "inst-1", "ubuntu:latest", "cid-1", string(domain.StatusRunning), "80:80", nil, nil, "10.0.0.1", "ovs-1", 1, now, now))

		inst, err := repo.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.NotNil(t, inst)
		assert.Equal(t, id, inst.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(id, userID).
			WillReturnError(pgx.ErrNoRows)

		inst, err := repo.GetByID(ctx, id)
		assert.Error(t, err)
		assert.Nil(t, inst)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestInstanceRepository_GetByName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		name := "inst-1"
		now := time.Now()

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(name, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "container_id", "status", "ports", "vpc_id", "subnet_id", "private_ip", "ovs_port", "version", "created_at", "updated_at"}).
				AddRow(id, userID, name, "ubuntu:latest", "cid-1", string(domain.StatusRunning), "80:80", nil, nil, "10.0.0.1", "ovs-1", 1, now, now))

		inst, err := repo.GetByName(ctx, name)
		assert.NoError(t, err)
		assert.NotNil(t, inst)
		assert.Equal(t, id, inst.ID)
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		name := "inst-1"

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(name, userID).
			WillReturnError(pgx.ErrNoRows)

		inst, err := repo.GetByName(ctx, name)
		assert.Error(t, err)
		assert.Nil(t, inst)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}

func TestInstanceRepository_List(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "container_id", "status", "ports", "vpc_id", "subnet_id", "private_ip", "ovs_port", "version", "created_at", "updated_at"}).
				AddRow(uuid.New(), userID, "inst-1", "ubuntu:latest", "cid-1", string(domain.StatusRunning), "80:80", nil, nil, "10.0.0.1", "ovs-1", 1, now, now))

		list, err := repo.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("db error", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(userID).
			WillReturnError(errors.New("db error"))

		list, err := repo.List(ctx)
		assert.Error(t, err)
		assert.Nil(t, list)
	})
}

func TestInstanceRepository_ListBySubnet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		userID := uuid.New()
		subnetID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)
		now := time.Now()

		mock.ExpectQuery("(?s)SELECT.+FROM instances.*").
			WithArgs(subnetID, userID).
			WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "image", "container_id", "status", "ports", "vpc_id", "subnet_id", "private_ip", "ovs_port", "version", "created_at", "updated_at"}).
				AddRow(uuid.New(), userID, "inst-1", "ubuntu:latest", "cid-1", string(domain.StatusRunning), "80:80", nil, nil, "10.0.0.1", "ovs-1", 1, now, now))

		list, err := repo.ListBySubnet(ctx, subnetID)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})
}

func TestInstanceRepository_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
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

		mock.ExpectExec("(?s)UPDATE instances.*").
			WithArgs(inst.Name, string(inst.Status), pgxmock.AnyArg(), inst.ContainerID, inst.Ports, inst.VpcID, inst.SubnetID, inst.PrivateIP, inst.OvsPort, inst.ID, inst.Version, inst.UserID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err = repo.Update(context.Background(), inst)
		assert.NoError(t, err)
		assert.Equal(t, 2, inst.Version)
	})

	t.Run("concurrency conflict", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		inst := &domain.Instance{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Status:  domain.StatusStopped,
			Version: 1,
		}

		mock.ExpectExec("(?s)UPDATE instances.*").
			WithArgs(inst.Name, string(inst.Status), pgxmock.AnyArg(), inst.ContainerID, inst.Ports, inst.VpcID, inst.SubnetID, inst.PrivateIP, inst.OvsPort, inst.ID, inst.Version, inst.UserID).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err = repo.Update(context.Background(), inst)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflict")
	})
}

func TestInstanceRepository_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectExec("DELETE FROM instances").
			WithArgs(id, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 1))

		err = repo.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mock, err := pgxmock.NewPool()
		assert.NoError(t, err)
		defer mock.Close()

		repo := NewInstanceRepository(mock)
		id := uuid.New()
		userID := uuid.New()
		ctx := appcontext.WithUserID(context.Background(), userID)

		mock.ExpectExec("DELETE FROM instances").
			WithArgs(id, userID).
			WillReturnResult(pgxmock.NewResult("DELETE", 0))

		err = repo.Delete(ctx, id)
		assert.Error(t, err)
		theCloudErr, ok := err.(*theclouderrors.Error)
		if ok {
			assert.Equal(t, theclouderrors.NotFound, theCloudErr.Type)
		}
	})
}
