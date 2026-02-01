package services_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/poyrazk/thecloud/internal/repositories/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const testRouteName = "test-route"

func TestGatewayServiceCreateRoute(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	cleanDB(t, db)

	repo := postgres.NewPostgresGatewayRepository(db)
	auditSvc := new(MockAuditService) // Keep audit mock for now as it probably writes to a different system or we don't care deeply about testing it right here
	auditSvc.On("Log", mock.Anything, mock.Anything, "gateway.route_create", "gateway", mock.Anything, mock.Anything).Return(nil)

	svc := services.NewGatewayService(repo, auditSvc)

	ctx := setupTestUser(t, db)

	params := ports.CreateRouteParams{
		Name:      testRouteName,
		Pattern:   "/test",
		Target:    "http://example.com",
		RateLimit: 100,
	}
	route, err := svc.CreateRoute(ctx, params)
	assert.NoError(t, err)
	assert.NotNil(t, route)
	assert.Equal(t, testRouteName, route.Name)
	assert.Equal(t, "prefix", route.PatternType) // Default prefix

	// Verify it exists in DB via repo
	userID := appcontext.UserIDFromContext(ctx)
	fetched, err := repo.GetRouteByID(ctx, route.ID, userID)
	assert.NoError(t, err)
	assert.Equal(t, route.ID, fetched.ID)
}

func TestGatewayServiceListRoutes(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	cleanDB(t, db)

	repo := postgres.NewPostgresGatewayRepository(db)
	auditSvc := new(MockAuditService)
	auditSvc.On("Log", mock.Anything, mock.Anything, "gateway.route_create", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := services.NewGatewayService(repo, auditSvc)
	ctx := setupTestUser(t, db)

	// Create a route
	params := ports.CreateRouteParams{Name: "r1", Pattern: "/r1", Target: "http://example.com"}
	_, err := svc.CreateRoute(ctx, params)
	require.NoError(t, err)

	res, err := svc.ListRoutes(ctx)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	assert.Equal(t, "r1", res[0].Name)
}

func TestGatewayServiceDeleteRoute(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	cleanDB(t, db)

	repo := postgres.NewPostgresGatewayRepository(db)
	auditSvc := new(MockAuditService)
	auditSvc.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := services.NewGatewayService(repo, auditSvc)
	ctx := setupTestUser(t, db)
	userID := appcontext.UserIDFromContext(ctx)

	// Create directly in repo or via service
	route := &domain.GatewayRoute{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "r1",
		PathPrefix:  "/r1",
		PathPattern: "/r1",
		PatternType: "prefix",
		TargetURL:   "http://example.com",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	err := repo.CreateRoute(ctx, route)
	require.NoError(t, err)

	// Refresh to sync cache (normally done by background worker or events, but here manual if needed?
	// CreateRoute updates cache, but manual repo insert doesn't.
	// But DeleteRoute calls repo.Delete then refreshes. So state in cache shouldn't matter for the deletion act itself,
	// except GetRouteByID check.

	err = svc.DeleteRoute(ctx, route.ID)
	assert.NoError(t, err)

	// Verify deleted
	_, err = repo.GetRouteByID(ctx, route.ID, userID)
	assert.Error(t, err)
}

func TestGatewayServiceGetProxy(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	cleanDB(t, db)

	repo := postgres.NewPostgresGatewayRepository(db)
	auditSvc := new(MockAuditService)
	auditSvc.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := services.NewGatewayService(repo, auditSvc)
	ctx := setupTestUser(t, db)

	// Create via service to ensure cache update
	params := ports.CreateRouteParams{Name: "api", Pattern: "/api", Target: "http://localhost:8080"}
	_, err := svc.CreateRoute(ctx, params)
	require.NoError(t, err)

	proxy, paramsMap, ok := svc.GetProxy("GET", "/api/users")
	assert.True(t, ok, "should match prefix /api")
	assert.NotNil(t, proxy)
	assert.Nil(t, paramsMap)

	_, _, ok = svc.GetProxy("GET", "/other")
	assert.False(t, ok)
}

func TestGatewayServiceGetProxyPattern(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	cleanDB(t, db)

	repo := postgres.NewPostgresGatewayRepository(db)
	auditSvc := new(MockAuditService)
	auditSvc.On("Log", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	svc := services.NewGatewayService(repo, auditSvc)
	ctx := setupTestUser(t, db)

	// Create pattern route via service
	params := ports.CreateRouteParams{Name: "users", Pattern: "/users/{id}", Target: "http://localhost:8080"}
	_, err := svc.CreateRoute(ctx, params)
	require.NoError(t, err)

	proxy, paramsMap, ok := svc.GetProxy("GET", "/users/123")
	assert.True(t, ok)
	assert.NotNil(t, proxy)
	assert.Equal(t, "123", paramsMap["id"])

	_, _, ok = svc.GetProxy("GET", "/users/123/posts")
	assert.False(t, ok, "should not match /users/123/posts with /users/{id} exactly unless wildcard used")
}
