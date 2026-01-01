package httphandlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poyraz/cloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type dashboardServiceMock struct {
	mock.Mock
}

func (m *dashboardServiceMock) GetSummary(ctx context.Context) (*domain.ResourceSummary, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ResourceSummary), args.Error(1)
}

func (m *dashboardServiceMock) GetRecentEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	args := m.Called(ctx, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Event), args.Error(1)
}

func (m *dashboardServiceMock) GetStats(ctx context.Context) (*domain.DashboardStats, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DashboardStats), args.Error(1)
}

// Actual test
func TestDashboardHandler_GetSummary(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(dashboardServiceMock)
	handler := NewDashboardHandler(mockSvc)

	r := gin.New()
	r.GET("/summary", handler.GetSummary)

	summary := &domain.ResourceSummary{
		TotalInstances:   5,
		RunningInstances: 3,
	}

	mockSvc.On("GetSummary", mock.Anything).Return(summary, nil)

	req, _ := http.NewRequest("GET", "/summary", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	// Handle httputil.Response wrapper
	var wrapper struct {
		Data domain.ResourceSummary `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &wrapper)
	assert.NoError(t, err)
	assert.Equal(t, 5, wrapper.Data.TotalInstances)
}

func TestDashboardHandler_GetRecentEvents(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockSvc := new(dashboardServiceMock)
	handler := NewDashboardHandler(mockSvc)

	r := gin.New()
	r.GET("/events", handler.GetRecentEvents)

	events := []*domain.Event{
		{ID: uuid.New(), Action: "TEST_ACTION"},
	}

	mockSvc.On("GetRecentEvents", mock.Anything, 10).Return(events, nil)

	req, _ := http.NewRequest("GET", "/events?limit=10", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var wrapper struct {
		Data []*domain.Event `json:"data"`
	}
	err := json.Unmarshal(w.Body.Bytes(), &wrapper)
	assert.NoError(t, err)
	assert.Len(t, wrapper.Data, 1)
	assert.Equal(t, "TEST_ACTION", wrapper.Data[0].Action)
}
