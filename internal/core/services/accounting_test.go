package services

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockAccountingRepository struct {
	mock.Mock
}

func (m *mockAccountingRepository) CreateRecord(ctx context.Context, record domain.UsageRecord) error {
	return m.Called(ctx, record).Error(0)
}

func (m *mockAccountingRepository) GetUsageSummary(ctx context.Context, userID uuid.UUID, start, end time.Time) (map[domain.ResourceType]float64, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).(map[domain.ResourceType]float64), args.Error(1)
}

func (m *mockAccountingRepository) ListRecords(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]domain.UsageRecord, error) {
	args := m.Called(ctx, userID, start, end)
	return args.Get(0).([]domain.UsageRecord), args.Error(1)
}

type mockInstanceRepository struct {
	mock.Mock
}

func (m *mockInstanceRepository) Create(ctx context.Context, instance *domain.Instance) error {
	return m.Called(ctx, instance).Error(0)
}
func (m *mockInstanceRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Instance, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepository) GetByName(ctx context.Context, name string) (*domain.Instance, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil { return nil, args.Error(1) }
	return args.Get(0).(*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepository) List(ctx context.Context) ([]*domain.Instance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepository) ListAll(ctx context.Context) ([]*domain.Instance, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepository) ListBySubnet(ctx context.Context, subnetID uuid.UUID) ([]*domain.Instance, error) {
    args := m.Called(ctx, subnetID)
    return args.Get(0).([]*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepository) Update(ctx context.Context, instance *domain.Instance) error {
	return m.Called(ctx, instance).Error(0)
}
func (m *mockInstanceRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

func TestAccountingService_TrackUsage(t *testing.T) {
	repo := new(mockAccountingRepository)
	svc := NewAccountingService(repo, nil, slog.Default())

	record := domain.UsageRecord{
		UserID:     uuid.New(),
		Quantity:   10,
		ResourceType: domain.ResourceInstance,
	}

	repo.On("CreateRecord", mock.Anything, mock.MatchedBy(func(r domain.UsageRecord) bool {
		return r.UserID == record.UserID && r.Quantity == record.Quantity
	})).Return(nil)

	err := svc.TrackUsage(context.Background(), record)
	assert.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestAccountingService_GetSummary(t *testing.T) {
	repo := new(mockAccountingRepository)
	svc := NewAccountingService(repo, nil, slog.Default())

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	usage := map[domain.ResourceType]float64{
		domain.ResourceInstance: 100,
		domain.ResourceStorage:  200,
	}

	repo.On("GetUsageSummary", mock.Anything, userID, start, end).Return(usage, nil)

	summary, err := svc.GetSummary(context.Background(), userID, start, end)
	assert.NoError(t, err)
	assert.NotNil(t, summary)
	// (100 * 0.01) + (200 * 0.005) = 1.0 + 1.0 = 2.0
	assert.Equal(t, 2.0, summary.TotalAmount)
	repo.AssertExpectations(t)
}

func TestAccountingService_ListUsage(t *testing.T) {
	repo := new(mockAccountingRepository)
	svc := NewAccountingService(repo, nil, slog.Default())

	userID := uuid.New()
	start := time.Now().Add(-24 * time.Hour)
	end := time.Now()

	records := []domain.UsageRecord{{ID: uuid.New()}}

	repo.On("ListRecords", mock.Anything, userID, start, end).Return(records, nil)

	res, err := svc.ListUsage(context.Background(), userID, start, end)
	assert.NoError(t, err)
	assert.Len(t, res, 1)
	repo.AssertExpectations(t)
}

func TestAccountingService_ProcessHourlyBilling(t *testing.T) {
	repo := new(mockAccountingRepository)
	instanceRepo := new(mockInstanceRepository)
	svc := NewAccountingService(repo, instanceRepo, slog.Default())

	instances := []*domain.Instance{
		{ID: uuid.New(), UserID: uuid.New(), Status: domain.StatusRunning},
		{ID: uuid.New(), UserID: uuid.New(), Status: domain.StatusStopped},
	}

	instanceRepo.On("ListAll", mock.Anything).Return(instances, nil)
	repo.On("CreateRecord", mock.Anything, mock.Anything).Return(nil)

	err := svc.ProcessHourlyBilling(context.Background())
	assert.NoError(t, err)
	
	// Should only record for running instances
	repo.AssertNumberOfCalls(t, "CreateRecord", 1)
	instanceRepo.AssertExpectations(t)
}
