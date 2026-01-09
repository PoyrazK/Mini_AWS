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

func TestVolumeRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	instanceID := uuid.New()
	vol := &domain.Volume{
		ID:         uuid.New(),
		UserID:     uuid.New(),
		Name:       "vol-1",
		SizeGB:     10,
		Status:     domain.VolumeStatusAvailable,
		InstanceID: &instanceID,
		MountPath:  "/mnt/data",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	mock.ExpectExec("INSERT INTO volumes").
		WithArgs(vol.ID, vol.UserID, vol.Name, vol.SizeGB, string(vol.Status), vol.InstanceID, vol.MountPath, vol.CreatedAt, vol.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), vol)
	assert.NoError(t, err)
}

func TestVolumeRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, size_gb, status, instance_id, mount_path, created_at, updated_at FROM volumes").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "size_gb", "status", "instance_id", "mount_path", "created_at", "updated_at"}).
			AddRow(id, userID, "vol-1", 10, string(domain.VolumeStatusAvailable), &id, "/mnt/data", now, now))

	vol, err := repo.GetByID(ctx, id)
	assert.NoError(t, err)
	assert.NotNil(t, vol)
	assert.Equal(t, id, vol.ID)
}

func TestVolumeRepository_GetByName(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()
	name := "vol-1"

	mock.ExpectQuery("SELECT id, user_id, name, size_gb, status, instance_id, mount_path, created_at, updated_at FROM volumes").
		WithArgs(name, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "size_gb", "status", "instance_id", "mount_path", "created_at", "updated_at"}).
			AddRow(id, userID, name, 10, string(domain.VolumeStatusAvailable), &id, "/mnt/data", now, now))

	vol, err := repo.GetByName(ctx, name)
	assert.NoError(t, err)
	assert.NotNil(t, vol)
	assert.Equal(t, id, vol.ID)
}

func TestVolumeRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	instID := uuid.New()
	mock.ExpectQuery("SELECT id, user_id, name, size_gb, status, instance_id, mount_path, created_at, updated_at FROM volumes").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "size_gb", "status", "instance_id", "mount_path", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "vol-1", 10, string(domain.VolumeStatusAvailable), &instID, "/mnt/data", now, now))

	vols, err := repo.List(ctx)
	assert.NoError(t, err)
	assert.Len(t, vols, 1)
}

func TestVolumeRepository_ListByInstanceID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	userID := uuid.New()
	instanceID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, size_gb, status, instance_id, mount_path, created_at, updated_at FROM volumes").
		WithArgs(instanceID, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "size_gb", "status", "instance_id", "mount_path", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "vol-1", 10, string(domain.VolumeStatusAvailable), &instanceID, "/mnt/data", now, now))

	vols, err := repo.ListByInstanceID(ctx, instanceID)
	assert.NoError(t, err)
	assert.Len(t, vols, 1)
}

func TestVolumeRepository_Update(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	vol := &domain.Volume{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Status:    domain.VolumeStatusInUse,
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec("UPDATE volumes").
		WithArgs(string(vol.Status), vol.InstanceID, vol.MountPath, vol.UpdatedAt, vol.ID, vol.UserID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.Update(context.Background(), vol)
	assert.NoError(t, err)
}

func TestVolumeRepository_Delete(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewVolumeRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)

	mock.ExpectExec("DELETE FROM volumes").
		WithArgs(id, userID).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err = repo.Delete(ctx, id)
	assert.NoError(t, err)
}
