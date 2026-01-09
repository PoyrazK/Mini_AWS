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

func TestGatewayRepository_CreateRoute(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresGatewayRepository(mock)
	route := &domain.GatewayRoute{
		ID:          uuid.New(),
		UserID:      uuid.New(),
		Name:        "test-route",
		PathPrefix:  "/api",
		TargetURL:   "http://localhost",
		StripPrefix: true,
		RateLimit:   100,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mock.ExpectExec("INSERT INTO gateway_routes").
		WithArgs(route.ID, route.UserID, route.Name, route.PathPrefix, route.TargetURL, route.StripPrefix, route.RateLimit, route.CreatedAt, route.UpdatedAt).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.CreateRoute(context.Background(), route)
	assert.NoError(t, err)
}

func TestGatewayRepository_GetRouteByID(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresGatewayRepository(mock)
	id := uuid.New()
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, path_prefix, target_url, strip_prefix, rate_limit, created_at, updated_at FROM gateway_routes").
		WithArgs(id, userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "path_prefix", "target_url", "strip_prefix", "rate_limit", "created_at", "updated_at"}).
			AddRow(id, userID, "test-route", "/api", "http://localhost", true, 100, now, now))

	route, err := repo.GetRouteByID(context.Background(), id, userID)
	assert.NoError(t, err)
	assert.NotNil(t, route)
	assert.Equal(t, id, route.ID)
}

func TestGatewayRepository_ListRoutes(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewPostgresGatewayRepository(mock)
	userID := uuid.New()
	now := time.Now()

	mock.ExpectQuery("SELECT id, user_id, name, path_prefix, target_url, strip_prefix, rate_limit, created_at, updated_at FROM gateway_routes").
		WithArgs(userID).
		WillReturnRows(pgxmock.NewRows([]string{"id", "user_id", "name", "path_prefix", "target_url", "strip_prefix", "rate_limit", "created_at", "updated_at"}).
			AddRow(uuid.New(), userID, "test-route", "/api", "http://localhost", true, 100, now, now))

	routes, err := repo.ListRoutes(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, routes, 1)
}
