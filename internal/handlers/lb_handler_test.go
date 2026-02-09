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
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	lbPath         = "/lb"
	testLBName     = "test-lb"
	algoRoundRobin = "round-robin"
)

type mockLBService struct {
	mock.Mock
}

func (m *mockLBService) Create(ctx context.Context, name string, vpcID uuid.UUID, port int, algorithm, idempotencyKey string) (*domain.LoadBalancer, error) {
	args := m.Called(ctx, name, vpcID, port, algorithm, idempotencyKey)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoadBalancer), args.Error(1)
}

func (m *mockLBService) List(ctx context.Context) ([]*domain.LoadBalancer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LoadBalancer), args.Error(1)
}

func (m *mockLBService) Get(ctx context.Context, id uuid.UUID) (*domain.LoadBalancer, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.LoadBalancer), args.Error(1)
}

func (m *mockLBService) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockLBService) AddTarget(ctx context.Context, lbID, instanceID uuid.UUID, port, weight int) error {
	args := m.Called(ctx, lbID, instanceID, port, weight)
	return args.Error(0)
}

func (m *mockLBService) RemoveTarget(ctx context.Context, lbID, instanceID uuid.UUID) error {
	args := m.Called(ctx, lbID, instanceID)
	return args.Error(0)
}

func (m *mockLBService) ListTargets(ctx context.Context, lbID uuid.UUID) ([]*domain.LBTarget, error) {
	args := m.Called(ctx, lbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LBTarget), args.Error(1)
}

func setupLBHandlerTest(_ *testing.T) (*mockLBService, *LBHandler, *gin.Engine) {
	gin.SetMode(gin.TestMode)
	svc := new(mockLBService)
	handler := NewLBHandler(svc)
	r := gin.New()
	return svc, handler, r
}

func TestLBHandlerCreate(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.POST(lbPath, handler.Create)

	vpcID := uuid.New()
	lb := &domain.LoadBalancer{ID: uuid.New(), Name: testLBName}
	svc.On("Create", mock.Anything, testLBName, vpcID, 80, algoRoundRobin, "").Return(lb, nil)

	body, err := json.Marshal(map[string]interface{}{
		"name":      testLBName,
		"vpc_id":    vpcID.String(),
		"port":      80,
		"algorithm": algoRoundRobin,
	})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", lbPath, bytes.NewBuffer(body))
	assert.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusAccepted, w.Code)
}

func TestLBHandlerList(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.GET(lbPath, handler.List)

	lbs := []*domain.LoadBalancer{{ID: uuid.New(), Name: "lb1"}}
	svc.On("List", mock.Anything).Return(lbs, nil)

	req := httptest.NewRequest(http.MethodGet, lbPath, nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLBHandlerGet(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.GET(lbPath+"/:id", handler.Get)

	id := uuid.New()
	lb := &domain.LoadBalancer{ID: id, Name: "lb1"}
	svc.On("Get", mock.Anything, id).Return(lb, nil)

	req := httptest.NewRequest(http.MethodGet, lbPath+"/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLBHandlerDelete(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.DELETE(lbPath+"/:id", handler.Delete)

	id := uuid.New()
	svc.On("Delete", mock.Anything, id).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, lbPath+"/"+id.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLBHandlerAddTarget(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.POST(lbPath+"/:id/targets", handler.AddTarget)

	lbID := uuid.New()
	instID := uuid.New()
	svc.On("AddTarget", mock.Anything, lbID, instID, 8080, 10).Return(nil)

	body, err := json.Marshal(map[string]interface{}{
		"instance_id": instID.String(),
		"port":        8080,
		"weight":      10,
	})
	assert.NoError(t, err)
	w := httptest.NewRecorder()
	req, err := http.NewRequest("POST", lbPath+"/"+lbID.String()+"/targets", bytes.NewBuffer(body))
	assert.NoError(t, err)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLBHandlerRemoveTarget(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.DELETE(lbPath+"/:id/targets/:instanceId", handler.RemoveTarget)

	lbID := uuid.New()
	instID := uuid.New()
	svc.On("RemoveTarget", mock.Anything, lbID, instID).Return(nil)

	req := httptest.NewRequest(http.MethodDelete, lbPath+"/"+lbID.String()+"/targets/"+instID.String(), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLBHandlerListTargets(t *testing.T) {
	t.Parallel()
	svc, handler, r := setupLBHandlerTest(t)
	defer svc.AssertExpectations(t)

	r.GET(lbPath+"/:id/targets", handler.ListTargets)

	lbID := uuid.New()
	targets := []*domain.LBTarget{{InstanceID: uuid.New()}}
	svc.On("ListTargets", mock.Anything, lbID).Return(targets, nil)

	req := httptest.NewRequest(http.MethodGet, lbPath+"/"+lbID.String()+"/targets", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestLBHandlerErrors(t *testing.T) {
	t.Parallel()

	t.Run("CreateInvalidJSON", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.POST(lbPath, handler.Create)
		req := httptest.NewRequest("POST", lbPath, bytes.NewBufferString("invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateInvalidVPC", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.POST(lbPath, handler.Create)
		body, _ := json.Marshal(map[string]interface{}{
			"name":   "n",
			"vpc_id": "invalid",
			"port":   80,
		})
		req := httptest.NewRequest("POST", lbPath, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("GetInvalidID", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.GET(lbPath+"/:id", handler.Get)
		req := httptest.NewRequest("GET", lbPath+"/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("DeleteInvalidID", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.DELETE(lbPath+"/:id", handler.Delete)
		req := httptest.NewRequest("DELETE", lbPath+"/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddTargetInvalidLBID", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.POST(lbPath+"/:id/targets", handler.AddTarget)
		req := httptest.NewRequest("POST", lbPath+"/invalid/targets", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("AddTargetInvalidJSON", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.POST(lbPath+"/:id/targets", handler.AddTarget)
		id := uuid.New()
		req := httptest.NewRequest("POST", lbPath+"/"+id.String()+"/targets", bytes.NewBufferString("invalid"))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("RemoveTargetInvalidID", func(t *testing.T) {
		_, handler, r := setupLBHandlerTest(t)
		r.DELETE(lbPath+"/:id/targets/:instanceId", handler.RemoveTarget)
		id := uuid.New()
		req := httptest.NewRequest("DELETE", lbPath+"/"+id.String()+"/targets/invalid", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("CreateServiceError", func(t *testing.T) {
		svc, handler, r := setupLBHandlerTest(t)
		r.POST(lbPath, handler.Create)
		svc.On("Create", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
			Return(nil, errors.New(errors.Internal, "internal error"))

		body, _ := json.Marshal(map[string]interface{}{
			"name":   "n",
			"vpc_id": uuid.New().String(),
			"port":   80,
		})
		req := httptest.NewRequest("POST", lbPath, bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("GetServiceError", func(t *testing.T) {
		svc, handler, r := setupLBHandlerTest(t)
		r.GET(lbPath+"/:id", handler.Get)
		id := uuid.New()
		svc.On("Get", mock.Anything, id).Return(nil, errors.New(errors.NotFound, "not found"))

		req := httptest.NewRequest("GET", lbPath+"/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusNotFound, w.Code)
	})

	t.Run("ListTargetsServiceError", func(t *testing.T) {
		svc, handler, r := setupLBHandlerTest(t)
		r.GET(lbPath+"/:id/targets", handler.ListTargets)
		id := uuid.New()
		svc.On("ListTargets", mock.Anything, id).Return(nil, errors.New(errors.Internal, "error"))

		req := httptest.NewRequest("GET", lbPath+"/"+id.String()+"/targets", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})

	t.Run("DeleteServiceError", func(t *testing.T) {
		svc, handler, r := setupLBHandlerTest(t)
		r.DELETE(lbPath+"/:id", handler.Delete)
		id := uuid.New()
		svc.On("Delete", mock.Anything, id).Return(errors.New(errors.Internal, "error"))

		req := httptest.NewRequest("DELETE", lbPath+"/"+id.String(), nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}
