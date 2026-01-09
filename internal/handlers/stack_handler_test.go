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

type mockStackService struct {
	mock.Mock
}

func (m *mockStackService) CreateStack(ctx context.Context, name, template string, parameters map[string]string) (*domain.Stack, error) {
	args := m.Called(ctx, name, template, parameters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Stack), args.Error(1)
}

func (m *mockStackService) ListStacks(ctx context.Context) ([]*domain.Stack, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Stack), args.Error(1)
}

func (m *mockStackService) GetStack(ctx context.Context, id uuid.UUID) (*domain.Stack, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Stack), args.Error(1)
}

func (m *mockStackService) DeleteStack(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockStackService) ValidateTemplate(ctx context.Context, template string) (*domain.TemplateValidateResponse, error) {
	args := m.Called(ctx, template)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.TemplateValidateResponse), args.Error(1)
}

func TestStackHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockStackService)
	handler := NewStackHandler(svc)

	stack := &domain.Stack{ID: uuid.New(), Name: "test-stack"}

	svc.On("CreateStack", mock.Anything, "test-stack", "version: 1", mock.Anything).Return(stack, nil)

	reqBody := CreateStackRequest{
		Name:     "test-stack",
		Template: "version: 1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/iac/stacks", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestStackHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockStackService)
	handler := NewStackHandler(svc)

	stacks := []*domain.Stack{
		{ID: uuid.New(), Name: "stack-1"},
	}

	svc.On("ListStacks", mock.Anything).Return(stacks, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/iac/stacks", nil)

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestStackHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockStackService)
	handler := NewStackHandler(svc)

	id := uuid.New()
	stack := &domain.Stack{ID: id, Name: "test-stack"}

	svc.On("GetStack", mock.Anything, id).Return(stack, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/iac/stacks/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestStackHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockStackService)
	handler := NewStackHandler(svc)

	id := uuid.New()

	svc.On("DeleteStack", mock.Anything, id).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/iac/stacks/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestStackHandler_Validate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockStackService)
	handler := NewStackHandler(svc)

	resp := &domain.TemplateValidateResponse{Valid: true}

	svc.On("ValidateTemplate", mock.Anything, "version: 1").Return(resp, nil)

	reqBody := map[string]string{
		"template": "version: 1",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/iac/validate", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Validate(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}
