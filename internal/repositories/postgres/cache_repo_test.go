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

func TestCacheRepository_Create(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCacheRepository(mock)
	cache := &domain.Cache{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "test-cache",
		Engine:      domain.EngineRedis,
		Version:     "6.0",
		Status:      domain.CacheStatusCreating,
		VpcID:       nil,
		ContainerID: "cid-1",
		Port:        6379,
		Password:    "password",
		MemoryMB:    1024,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO caches").
		WithArgs(cache.ID, cache.UserID, cache.Name, cache.Engine, cache.Version, cache.Status, cache.VpcID,
			cache.ContainerID, cache.Port, cache.Password, cache.MemoryMB, cache.CreatedAt, cache.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), cache)
	assert.NoError(t, err)
}

func TestCacheRepository_GetByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCacheRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, engine, version, status, vpc_id, container_id, port, password, memory_mb, created_at, updated_at FROM caches").
		WithArgs(id).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "engine", "version", "status", "vpc_id", "container_id", "port", "password", "memory_mb", "created_at", "updated_at"}).
			AddRow(id, userID, "test-cache", string(domain.EngineRedis), "6.0", string(domain.CacheStatusCreating), nil, "cid-1", 6379, "password", 1024, now, now))

	c, err := repo.GetByID(context.Background(), id)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.Equal(t, id, c.ID)
}

func TestCacheRepository_List(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewCacheRepository(mock)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, engine, version, status, vpc_id, container_id, port, password, memory_mb, created_at, updated_at FROM caches").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "engine", "version", "status", "vpc_id", "container_id", "port", "password", "memory_mb", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "test-cache", string(domain.EngineRedis), "6.0", string(domain.CacheStatusCreating), nil, "cid-1", 6379, "password", 1024, now, now))

	caches, err := repo.List(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, caches, 1)
}
