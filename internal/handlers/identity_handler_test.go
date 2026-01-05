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

type mockIdentityService struct {
	mock.Mock
}

func (m *mockIdentityService) CreateKey(ctx context.Context, userID uuid.UUID, name string) (*domain.APIKey, error) {
	args := m.Called(ctx, userID, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func (m *mockIdentityService) ValidateAPIKey(ctx context.Context, key string) (*domain.APIKey, error) {
	args := m.Called(ctx, key)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}
func (m *mockIdentityService) ListKeys(ctx context.Context, userID uuid.UUID) ([]*domain.APIKey, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.APIKey), args.Error(1)
}
func (m *mockIdentityService) RevokeKey(ctx context.Context, userID, id uuid.UUID) error {
	args := m.Called(ctx, userID, id)
	return args.Error(0)
}
func (m *mockIdentityService) RotateKey(ctx context.Context, userID, id uuid.UUID) (*domain.APIKey, error) {
	args := m.Called(ctx, userID, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.APIKey), args.Error(1)
}

func TestIdentityHandler_CreateKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockIdentityService)
	handler := NewIdentityHandler(svc)

	userID := uuid.New()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.POST("/auth/keys", handler.CreateKey)

	key := &domain.APIKey{Key: "sk_test_123", Name: "Test Key"}
	svc.On("CreateKey", mock.Anything, userID, "Test Key").Return(key, nil)

	body, _ := json.Marshal(map[string]string{"name": "Test Key"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/keys", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "sk_test_123")
}

func TestIdentityHandler_ListKeys(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockIdentityService)
	handler := NewIdentityHandler(svc)

	userID := uuid.New()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.GET("/auth/keys", handler.ListKeys)

	keys := []*domain.APIKey{{ID: uuid.New(), Name: "K1"}, {ID: uuid.New(), Name: "K2"}}
	svc.On("ListKeys", mock.Anything, userID).Return(keys, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/auth/keys", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "K1")
	assert.Contains(t, w.Body.String(), "K2")
}

func TestIdentityHandler_RevokeKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockIdentityService)
	handler := NewIdentityHandler(svc)

	userID := uuid.New()
	keyID := uuid.New()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("userID", userID)
		c.Next()
	})
	r.DELETE("/auth/keys/:id", handler.RevokeKey)

	svc.On("RevokeKey", mock.Anything, userID, keyID).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("DELETE", "/auth/keys/"+keyID.String(), nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
}

func TestIdentityHandler_RotateKey(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockIdentityService)
	handler := NewIdentityHandler(svc)

	userID := uuid.New()
	keyID := uuid.New()
	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	r.POST("/auth/keys/:id/rotate", handler.RotateKey)

	newKey := &domain.APIKey{ID: uuid.New(), Key: "new_key", Name: "Rotated"}
	svc.On("RotateKey", mock.Anything, mock.Anything, keyID).Return(newKey, nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/auth/keys/"+keyID.String()+"/rotate", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "new_key")
}
