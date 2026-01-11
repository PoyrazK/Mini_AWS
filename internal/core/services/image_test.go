package services_test

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockImageRepo struct {
	mock.Mock
}

func (m *mockImageRepo) Create(ctx context.Context, img *domain.Image) error {
	args := m.Called(ctx, img)
	return args.Error(0)
}

func (m *mockImageRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Image, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Image), args.Error(1)
}

func (m *mockImageRepo) List(ctx context.Context, userID uuid.UUID, includePublic bool) ([]*domain.Image, error) {
	args := m.Called(ctx, userID, includePublic)
	return args.Get(0).([]*domain.Image), args.Error(1)
}

func (m *mockImageRepo) Update(ctx context.Context, img *domain.Image) error {
	args := m.Called(ctx, img)
	return args.Error(0)
}

func (m *mockImageRepo) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type mockFileStore struct {
	mock.Mock
}

func (m *mockFileStore) Write(ctx context.Context, bucket, key string, r io.Reader) (int64, error) {
	args := m.Called(ctx, bucket, key, r)
	return int64(args.Int(0)), args.Error(1)
}

func (m *mockFileStore) Read(ctx context.Context, bucket, key string) (io.ReadCloser, error) {
	args := m.Called(ctx, bucket, key)
	return args.Get(0).(io.ReadCloser), args.Error(1)
}

func (m *mockFileStore) Delete(ctx context.Context, bucket, key string) error {
	args := m.Called(ctx, bucket, key)
	return args.Error(0)
}

func TestImageService(t *testing.T) {
	repo := new(mockImageRepo)
	store := new(mockFileStore)
	svc := services.NewImageService(repo, store, nil)
	userID := uuid.New()
	ctx := appcontext.WithUserID(context.Background(), userID)

	t.Run("RegisterImage", func(t *testing.T) {
		repo.On("Create", mock.Anything, mock.Anything).Return(nil)
		img, err := svc.RegisterImage(ctx, "ubuntu", "Ubuntu 22.04", "linux", "22.04", true)
		assert.NoError(t, err)
		assert.NotNil(t, img)
		assert.Equal(t, "ubuntu", img.Name)
		repo.AssertExpectations(t)
	})

	t.Run("UploadImage", func(t *testing.T) {
		id := uuid.New()
		img := &domain.Image{ID: id, UserID: userID}
		repo.On("GetByID", mock.Anything, id).Return(img, nil)
		store.On("Write", mock.Anything, "images", mock.Anything, mock.Anything).Return(1024, nil)
		repo.On("Update", mock.Anything, img).Return(nil)

		err := svc.UploadImage(ctx, id, strings.NewReader("fake content"))
		assert.NoError(t, err)
		assert.Equal(t, domain.ImageStatusActive, img.Status)
	})
}
