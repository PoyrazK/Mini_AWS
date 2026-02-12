package workers

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/stretchr/testify/mock"
)

type mockInstanceRepo struct {
	mock.Mock
}

func (m *mockInstanceRepo) Create(ctx context.Context, instance *domain.Instance) error { return nil }
func (m *mockInstanceRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceRepo) GetByName(ctx context.Context, name string) (*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceRepo) List(ctx context.Context) ([]*domain.Instance, error) { return nil, nil }
func (m *mockInstanceRepo) ListAll(ctx context.Context) ([]*domain.Instance, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Instance), args.Error(1)
}
func (m *mockInstanceRepo) ListBySubnet(ctx context.Context, subnetID uuid.UUID) ([]*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceRepo) Update(ctx context.Context, instance *domain.Instance) error { return nil }
func (m *mockInstanceRepo) Delete(ctx context.Context, id uuid.UUID) error              { return nil }

type mockInstanceSvc struct {
	mock.Mock
}

func (m *mockInstanceSvc) LaunchInstance(ctx context.Context, params ports.LaunchParams) (*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceSvc) LaunchInstanceWithOptions(ctx context.Context, opts ports.CreateInstanceOptions) (*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceSvc) StartInstance(ctx context.Context, idOrName string) error {
	return m.Called(ctx, idOrName).Error(0)
}
func (m *mockInstanceSvc) StopInstance(ctx context.Context, idOrName string) error {
	return m.Called(ctx, idOrName).Error(0)
}
func (m *mockInstanceSvc) ListInstances(ctx context.Context) ([]*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceSvc) GetInstance(ctx context.Context, idOrName string) (*domain.Instance, error) {
	return nil, nil
}
func (m *mockInstanceSvc) GetInstanceLogs(ctx context.Context, idOrName string) (string, error) {
	return "", nil
}
func (m *mockInstanceSvc) GetInstanceStats(ctx context.Context, idOrName string) (*domain.InstanceStats, error) {
	return nil, nil
}
func (m *mockInstanceSvc) GetConsoleURL(ctx context.Context, idOrName string) (string, error) {
	return "", nil
}
func (m *mockInstanceSvc) TerminateInstance(ctx context.Context, idOrName string) error { return nil }
func (m *mockInstanceSvc) Exec(ctx context.Context, idOrName string, cmd []string) (string, error) {
	return "", nil
}
func (m *mockInstanceSvc) UpdateInstanceMetadata(ctx context.Context, id uuid.UUID, metadata, labels map[string]string) error {
	return nil
}

func TestHealingWorker(t *testing.T) {
	repo := new(mockInstanceRepo)
	svc := new(mockInstanceSvc)
	logger := slog.Default()

	worker := NewHealingWorker(svc, repo, logger)
	worker.healingDelay = 1 * time.Millisecond

	t.Run("Heal Error Instances", func(t *testing.T) {
		inst1 := &domain.Instance{ID: uuid.New(), Status: domain.StatusRunning}
		inst2 := &domain.Instance{ID: uuid.New(), Status: domain.StatusError, UserID: uuid.New(), TenantID: uuid.New()}

		repo.On("ListAll", mock.Anything).Return([]*domain.Instance{inst1, inst2}, nil)
		svc.On("StopInstance", mock.Anything, inst2.ID.String()).Return(nil)
		svc.On("StartInstance", mock.Anything, inst2.ID.String()).Return(nil)

		worker.healERRORInstances(context.Background())

		// Wait for async healing tasks to complete
		worker.reconcileWG.Wait()

		repo.AssertExpectations(t)
		svc.AssertExpectations(t)

		svc.AssertNotCalled(t, "StopInstance", mock.Anything, inst1.ID.String())
	})
}
