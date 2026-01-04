package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGatewayService_CreateRoute(t *testing.T) {
	repo := new(MockGatewayRepo)
	repo.On("GetAllActiveRoutes", mock.Anything).Return([]*domain.GatewayRoute{}, nil)
	svc := services.NewGatewayService(repo)

	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)

	repo.On("CreateRoute", ctx, mock.AnythingOfType("*domain.GatewayRoute")).Return(nil)
	repo.On("GetAllActiveRoutes", ctx).Return([]*domain.GatewayRoute{}, nil)

	route, err := svc.CreateRoute(ctx, "test-api", "/v1", "http://target:80", true, 100)

	assert.NoError(t, err)
	assert.NotNil(t, route)
	assert.Equal(t, "/v1", route.PathPrefix)
	repo.AssertExpectations(t)
}

func TestGatewayService_RefreshAndGetProxy(t *testing.T) {
	repo := new(MockGatewayRepo)

	route := &domain.GatewayRoute{
		PathPrefix: "/api",
		TargetURL:  "http://localhost:8080",
	}
	repo.On("GetAllActiveRoutes", mock.Anything).Return([]*domain.GatewayRoute{route}, nil)

	svc := services.NewGatewayService(repo)

	proxy, ok := svc.GetProxy("/api/users")
	assert.True(t, ok)
	assert.NotNil(t, proxy)

	_, ok = svc.GetProxy("/wrong")
	assert.False(t, ok)
}
