package httphandlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockElasticIPService struct {
	mock.Mock
}

func (m *mockElasticIPService) AllocateIP(ctx context.Context) (*domain.ElasticIP, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ElasticIP), args.Error(1)
}

func (m *mockElasticIPService) ReleaseIP(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockElasticIPService) AssociateIP(ctx context.Context, id uuid.UUID, instanceID uuid.UUID) (*domain.ElasticIP, error) {
	args := m.Called(ctx, id, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ElasticIP), args.Error(1)
}

func (m *mockElasticIPService) DisassociateIP(ctx context.Context, id uuid.UUID) (*domain.ElasticIP, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ElasticIP), args.Error(1)
}

func (m *mockElasticIPService) ListElasticIPs(ctx context.Context) ([]*domain.ElasticIP, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.ElasticIP), args.Error(1)
}

func (m *mockElasticIPService) GetElasticIP(ctx context.Context, id uuid.UUID) (*domain.ElasticIP, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.ElasticIP), args.Error(1)
}

func setupElasticIPHandlerTest() (*mockElasticIPService, *ElasticIPHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	svc := new(mockElasticIPService)
	handler := NewElasticIPHandler(svc)
	r := gin.New()
	return svc, handler, r
}

func TestElasticIPHandlerAllocate(t *testing.T) {
	svc, handler, r := setupElasticIPHandlerTest()
	r.POST("/elastic-ips", handler.Allocate)

	eip := &domain.ElasticIP{ID: uuid.New(), PublicIP: "1.2.3.4"}
	svc.On("AllocateIP", mock.Anything).Return(eip, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/elastic-ips", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestElasticIPHandlerList(t *testing.T) {
	svc, handler, r := setupElasticIPHandlerTest()
	r.GET("/elastic-ips", handler.List)

	svc.On("ListElasticIPs", mock.Anything).Return([]*domain.ElasticIP{}, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/elastic-ips", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestElasticIPHandlerAssociate(t *testing.T) {
	svc, handler, r := setupElasticIPHandlerTest()
	r.POST("/elastic-ips/:id/associate", handler.Associate)

	eipID := uuid.New()
	instID := uuid.New()
	svc.On("AssociateIP", mock.Anything, eipID, instID).Return(&domain.ElasticIP{}, nil)

	body, _ := json.Marshal(map[string]string{"instance_id": instID.String()})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/elastic-ips/"+eipID.String()+"/associate", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
