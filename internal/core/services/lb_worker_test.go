package services

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockLBRepo struct {
	mock.Mock
}

func (m *mockLBRepo) Create(ctx context.Context, lb *domain.LoadBalancer) error { return nil }
func (m *mockLBRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.LoadBalancer, error) {
	return nil, nil
}
func (m *mockLBRepo) GetByIdempotencyKey(ctx context.Context, key string) (*domain.LoadBalancer, error) {
	return nil, nil
}
func (m *mockLBRepo) List(ctx context.Context) ([]*domain.LoadBalancer, error) { return nil, nil }
func (m *mockLBRepo) ListAll(ctx context.Context) ([]*domain.LoadBalancer, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LoadBalancer), args.Error(1)
}
func (m *mockLBRepo) Update(ctx context.Context, lb *domain.LoadBalancer) error {
	return m.Called(ctx, lb).Error(0)
}
func (m *mockLBRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *mockLBRepo) AddTarget(ctx context.Context, target *domain.LBTarget) error       { return nil }
func (m *mockLBRepo) RemoveTarget(ctx context.Context, lbID, instanceID uuid.UUID) error { return nil }
func (m *mockLBRepo) ListTargets(ctx context.Context, lbID uuid.UUID) ([]*domain.LBTarget, error) {
	args := m.Called(ctx, lbID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.LBTarget), args.Error(1)
}
func (m *mockLBRepo) UpdateTargetHealth(ctx context.Context, lbID, instanceID uuid.UUID, health string) error {
	return m.Called(ctx, lbID, instanceID, health).Error(0)
}
func (m *mockLBRepo) GetTargetsForInstance(ctx context.Context, instanceID uuid.UUID) ([]*domain.LBTarget, error) {
	return nil, nil
}

type mockProxyAdapter struct {
	mock.Mock
}

func (m *mockProxyAdapter) DeployProxy(ctx context.Context, lb *domain.LoadBalancer, targets []*domain.LBTarget) (string, error) {
	args := m.Called(ctx, lb, targets)
	return args.String(0), args.Error(1)
}

func (m *mockProxyAdapter) RemoveProxy(ctx context.Context, lbID uuid.UUID) error {
	return m.Called(ctx, lbID).Error(0)
}

func (m *mockProxyAdapter) UpdateProxyConfig(ctx context.Context, lb *domain.LoadBalancer, targets []*domain.LBTarget) error {
	return m.Called(ctx, lb, targets).Error(0)
}

func TestLBWorker_ProcessCreatingLBs(t *testing.T) {
	lbRepo := new(mockLBRepo)
	proxyAdapter := new(mockProxyAdapter)
	worker := NewLBWorker(lbRepo, nil, proxyAdapter)
	ctx := context.Background()
	lbID := uuid.New()
	userID := uuid.New()
	lb := &domain.LoadBalancer{ID: lbID, UserID: userID, Status: domain.LBStatusCreating}
	targets := []*domain.LBTarget{{InstanceID: uuid.New(), Port: 80}}

	lbRepo.On("ListAll", ctx).Return([]*domain.LoadBalancer{lb}, nil)
	lbRepo.On("ListTargets", mock.Anything, lbID).Return(targets, nil)
	proxyAdapter.On("DeployProxy", mock.Anything, lb, targets).Return("container-123", nil)
	lbRepo.On("Update", mock.Anything, mock.MatchedBy(func(l *domain.LoadBalancer) bool {
		return l.Status == domain.LBStatusActive
	})).Return(nil)

	worker.processCreatingLBs(ctx)
	lbRepo.AssertExpectations(t)
	proxyAdapter.AssertExpectations(t)
}

func TestLBWorker_ProcessDeletingLBs(t *testing.T) {
	lbRepo := new(mockLBRepo)
	proxyAdapter := new(mockProxyAdapter)
	worker := NewLBWorker(lbRepo, nil, proxyAdapter)
	ctx := context.Background()
	lbID := uuid.New()
	userID := uuid.New()
	lb := &domain.LoadBalancer{ID: lbID, UserID: userID, Status: domain.LBStatusDeleted}

	lbRepo.On("ListAll", ctx).Return([]*domain.LoadBalancer{lb}, nil)
	proxyAdapter.On("RemoveProxy", mock.Anything, lbID).Return(nil)
	lbRepo.On("Delete", mock.Anything, lbID).Return(nil)

	worker.processDeletingLBs(ctx)
	lbRepo.AssertExpectations(t)
	proxyAdapter.AssertExpectations(t)
}

func TestLBWorker_ProcessActiveLBs(t *testing.T) {
	lbRepo := new(mockLBRepo)
	proxyAdapter := new(mockProxyAdapter)
	worker := NewLBWorker(lbRepo, nil, proxyAdapter)
	ctx := context.Background()
	lbID := uuid.New()
	userID := uuid.New()
	lb := &domain.LoadBalancer{ID: lbID, UserID: userID, Status: domain.LBStatusActive}
	targets := []*domain.LBTarget{{InstanceID: uuid.New(), Port: 80}}

	lbRepo.On("ListAll", ctx).Return([]*domain.LoadBalancer{lb}, nil)
	lbRepo.On("ListTargets", mock.Anything, lbID).Return(targets, nil)
	proxyAdapter.On("UpdateProxyConfig", mock.Anything, lb, targets).Return(nil)

	worker.processActiveLBs(ctx)
	lbRepo.AssertExpectations(t)
	proxyAdapter.AssertExpectations(t)
}

func TestLBWorker_ProcessHealthChecks(t *testing.T) {
	lbRepo := new(mockLBRepo)
	instRepo := new(mockInstanceRepo)
	worker := NewLBWorker(lbRepo, instRepo, nil)
	ctx := context.Background()
	lbID := uuid.New()
	userID := uuid.New()
	lb := &domain.LoadBalancer{ID: lbID, UserID: userID, Status: domain.LBStatusActive}
	instID := uuid.New()
	targets := []*domain.LBTarget{{InstanceID: instID, Port: 80, Health: "unknown"}}
	inst := &domain.Instance{ID: instID, Ports: "49151:80"}

	lbRepo.On("ListAll", ctx).Return([]*domain.LoadBalancer{lb}, nil)
	lbRepo.On("ListTargets", mock.Anything, lbID).Return(targets, nil)
	instRepo.On("GetByID", mock.Anything, instID).Return(inst, nil)

	lbRepo.On("UpdateTargetHealth", mock.Anything, lbID, instID, mock.Anything).Return(nil)

	worker.processHealthChecks(ctx)
	lbRepo.AssertExpectations(t)
}

func TestLBWorker_Errors(t *testing.T) {
	ctx := context.Background()

	t.Run("ProcessCreating_ListError", func(t *testing.T) {
		lbRepo := new(mockLBRepo)
		worker := NewLBWorker(lbRepo, nil, nil)
		lbRepo.On("ListAll", ctx).Return(nil, assert.AnError)
		worker.processCreatingLBs(ctx)
		lbRepo.AssertExpectations(t)
	})

	t.Run("DeployLB_ListTargetsError", func(t *testing.T) {
		lbRepo := new(mockLBRepo)
		worker := NewLBWorker(lbRepo, nil, nil)
		lb := &domain.LoadBalancer{ID: uuid.New()}
		lbRepo.On("ListTargets", mock.Anything, lb.ID).Return(nil, assert.AnError)
		worker.deployLB(ctx, lb)
		lbRepo.AssertExpectations(t)
	})

	t.Run("DeployLB_DeployError", func(t *testing.T) {
		lbRepo := new(mockLBRepo)
		proxyAdapter := new(mockProxyAdapter)
		worker := NewLBWorker(lbRepo, nil, proxyAdapter)
		lb := &domain.LoadBalancer{ID: uuid.New()}
		lbRepo.On("ListTargets", mock.Anything, lb.ID).Return([]*domain.LBTarget{}, nil)
		proxyAdapter.On("DeployProxy", mock.Anything, lb, mock.Anything).Return("", assert.AnError)
		worker.deployLB(ctx, lb)
	})
}
