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

type mockSnapshotService struct {
	mock.Mock
}

func (m *mockSnapshotService) CreateSnapshot(ctx context.Context, volumeID uuid.UUID, description string) (*domain.Snapshot, error) {
	args := m.Called(ctx, volumeID, description)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Snapshot), args.Error(1)
}

func (m *mockSnapshotService) ListSnapshots(ctx context.Context) ([]*domain.Snapshot, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Snapshot), args.Error(1)
}

func (m *mockSnapshotService) GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Snapshot), args.Error(1)
}

func (m *mockSnapshotService) DeleteSnapshot(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockSnapshotService) RestoreSnapshot(ctx context.Context, id uuid.UUID, newVolumeName string) (*domain.Volume, error) {
	args := m.Called(ctx, id, newVolumeName)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Volume), args.Error(1)
}

func TestSnapshotHandler_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockSnapshotService)
	handler := NewSnapshotHandler(svc)

	volumeID := uuid.New()
	snapshot := &domain.Snapshot{ID: uuid.New(), VolumeID: volumeID, Description: "test snapshot"}

	svc.On("CreateSnapshot", mock.Anything, volumeID, "test snapshot").Return(snapshot, nil)

	reqBody := CreateSnapshotRequest{
		VolumeID:    volumeID,
		Description: "test snapshot",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/snapshots", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Create(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}

func TestSnapshotHandler_List(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockSnapshotService)
	handler := NewSnapshotHandler(svc)

	snapshots := []*domain.Snapshot{
		{ID: uuid.New()},
		{ID: uuid.New()},
	}

	svc.On("ListSnapshots", mock.Anything).Return(snapshots, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/snapshots", nil)

	handler.List(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestSnapshotHandler_Get(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockSnapshotService)
	handler := NewSnapshotHandler(svc)

	id := uuid.New()
	snapshot := &domain.Snapshot{ID: id}

	svc.On("GetSnapshot", mock.Anything, id).Return(snapshot, nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/snapshots/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Get(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestSnapshotHandler_Delete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockSnapshotService)
	handler := NewSnapshotHandler(svc)

	id := uuid.New()

	svc.On("DeleteSnapshot", mock.Anything, id).Return(nil)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("DELETE", "/snapshots/"+id.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Delete(c)

	assert.Equal(t, http.StatusOK, w.Code)
	svc.AssertExpectations(t)
}

func TestSnapshotHandler_Restore(t *testing.T) {
	gin.SetMode(gin.TestMode)
	svc := new(mockSnapshotService)
	handler := NewSnapshotHandler(svc)

	id := uuid.New()
	vol := &domain.Volume{ID: uuid.New(), Name: "restored-volume"}

	svc.On("RestoreSnapshot", mock.Anything, id, "restored-volume").Return(vol, nil)

	reqBody := RestoreSnapshotRequest{
		NewVolumeName: "restored-volume",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/snapshots/"+id.String()+"/restore", bytes.NewBuffer(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: id.String()}}

	handler.Restore(c)

	assert.Equal(t, http.StatusCreated, w.Code)
	svc.AssertExpectations(t)
}
