package services

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockDatabaseRepo struct {
	mock.Mock
}

func (m *mockDatabaseRepo) Create(ctx context.Context, db *domain.Database) error {
	return m.Called(ctx, db).Error(0)
}
func (m *mockDatabaseRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Database, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Database), args.Error(1)
}
func (m *mockDatabaseRepo) List(ctx context.Context) ([]*domain.Database, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Database), args.Error(1)
}
func (m *mockDatabaseRepo) ListReplicas(ctx context.Context, primaryID uuid.UUID) ([]*domain.Database, error) {
	args := m.Called(ctx, primaryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Database), args.Error(1)
}
func (m *mockDatabaseRepo) Update(ctx context.Context, db *domain.Database) error {
	return m.Called(ctx, db).Error(0)
}
func (m *mockDatabaseRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}

type mockCompute struct {
	mock.Mock
	ports.ComputeBackend
}

func (m *mockCompute) LaunchInstanceWithOptions(ctx context.Context, opts ports.CreateInstanceOptions) (string, []string, error) {
	args := m.Called(ctx, opts)
	return args.String(0), args.Get(1).([]string), args.Error(2)
}
func (m *mockCompute) GetInstanceIP(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}
func (m *mockCompute) GetInstancePort(ctx context.Context, id string, port string) (int, error) {
	args := m.Called(ctx, id, port)
	return args.Int(0), args.Error(1)
}
func (m *mockCompute) DeleteInstance(ctx context.Context, id string) error {
	return m.Called(ctx, id).Error(0)
}

type mockDatabaseVpcRepo struct {
	mock.Mock
}

func (m *mockDatabaseVpcRepo) Create(ctx context.Context, vpc *domain.VPC) error { return nil }
func (m *mockDatabaseVpcRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.VPC, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.VPC), args.Error(1)
}
func (m *mockDatabaseVpcRepo) GetByName(ctx context.Context, name string) (*domain.VPC, error) { return nil, nil }
func (m *mockDatabaseVpcRepo) List(ctx context.Context) ([]*domain.VPC, error)                { return nil, nil }
func (m *mockDatabaseVpcRepo) Delete(ctx context.Context, id uuid.UUID) error                 { return nil }
func (m *mockDatabaseVpcRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.VPC, error) {
	return nil, nil
}

type mockEventSvc struct {
	mock.Mock
}

func (m *mockEventSvc) RecordEvent(ctx context.Context, action, resourceID, resourceType string, metadata map[string]interface{}) error {
	return nil
}
func (m *mockEventSvc) ListEvents(ctx context.Context, limit int) ([]*domain.Event, error) {
	return nil, nil
}

type mockAuditSvc struct {
	mock.Mock
}

func (m *mockAuditSvc) Log(ctx context.Context, userID uuid.UUID, action, resourceType, resourceID string, details map[string]interface{}) error {
	return nil
}
func (m *mockAuditSvc) ListLogs(ctx context.Context, userID uuid.UUID, limit int) ([]*domain.AuditLog, error) {
	return nil, nil
}

func TestDatabaseService_Unit(t *testing.T) {
	repo := new(mockDatabaseRepo)
	compute := new(mockCompute)
	vpcRepo := new(mockDatabaseVpcRepo)
	eventSvc := new(mockEventSvc)
	auditSvc := new(mockAuditSvc)
	logger := slog.Default()

	svc := NewDatabaseService(DatabaseServiceParams{
		Repo:     repo,
		Compute:  compute,
		VpcRepo:  vpcRepo,
		EventSvc: eventSvc,
		AuditSvc: auditSvc,
		Logger:   logger,
	})

	ctx := context.Background()

	t.Run("CreateReplica", func(t *testing.T) {
		primaryID := uuid.New()
		primary := &domain.Database{
			ID:          primaryID,
			Engine:      domain.EnginePostgres,
			Version:     "16",
			Username:    "user",
			Password:    "pass",
			ContainerID: "cid-primary",
		}

		repo.On("GetByID", mock.Anything, primaryID).Return(primary, nil)
		compute.On("GetInstanceIP", mock.Anything, "cid-primary").Return("10.0.0.1", nil)
		compute.On("LaunchInstanceWithOptions", mock.Anything, mock.Anything).Return("cid-replica", []string{"5432:5432"}, nil)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(db *domain.Database) bool {
			return db.Role == domain.RoleReplica && *db.PrimaryID == primaryID
		})).Return(nil)

		db, err := svc.CreateReplica(ctx, primaryID, "my-replica")
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.Equal(t, domain.RoleReplica, db.Role)

		repo.AssertExpectations(t)
		compute.AssertExpectations(t)
	})

	t.Run("PromoteToPrimary", func(t *testing.T) {
		id := uuid.New()
		replica := &domain.Database{
			ID:   id,
			Role: domain.RoleReplica,
		}

		repo.On("GetByID", mock.Anything, id).Return(replica, nil)
		repo.On("Update", mock.Anything, mock.MatchedBy(func(db *domain.Database) bool {
			return db.Role == domain.RolePrimary && db.PrimaryID == nil
		})).Return(nil)

		err := svc.PromoteToPrimary(ctx, id)
		assert.NoError(t, err)

		repo.AssertExpectations(t)
	})

	t.Run("CreateDatabase", func(t *testing.T) {
		compute.On("LaunchInstanceWithOptions", mock.Anything, mock.Anything).Return("cid-1", []string{"5432:5432"}, nil)
		repo.On("Create", mock.Anything, mock.MatchedBy(func(db *domain.Database) bool {
			return db.Name == "new-db" && db.Role == domain.RolePrimary
		})).Return(nil)

		db, err := svc.CreateDatabase(ctx, "new-db", "postgres", "16", nil)
		assert.NoError(t, err)
		assert.NotNil(t, db)
		assert.Equal(t, domain.RolePrimary, db.Role)

		repo.AssertExpectations(t)
	})

	t.Run("GetAndList", func(t *testing.T) {
		id := uuid.New()
		db := &domain.Database{ID: id}
		repo.On("GetByID", mock.Anything, id).Return(db, nil)
		repo.On("List", mock.Anything).Return([]*domain.Database{db}, nil)

		res, err := svc.GetDatabase(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)

		list, err := svc.ListDatabases(ctx)
		assert.NoError(t, err)
		assert.Len(t, list, 1)
	})

	t.Run("DeleteDatabase", func(t *testing.T) {
		id := uuid.New()
		db := &domain.Database{ID: id, UserID: uuid.New(), ContainerID: "cid-1"}
		repo.On("GetByID", mock.Anything, id).Return(db, nil)
		compute.On("DeleteInstance", mock.Anything, "cid-1").Return(nil)
		repo.On("Delete", mock.Anything, id).Return(nil)

		err := svc.DeleteDatabase(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("GetConnectionString", func(t *testing.T) {
		id := uuid.New()
		db := &domain.Database{
			ID:       id,
			Engine:   domain.EnginePostgres,
			Username: "u",
			Password: "p",
			Port:     5432,
			Name:     "d",
		}
		repo.On("GetByID", mock.Anything, id).Return(db, nil)

		conn, err := svc.GetConnectionString(ctx, id)
		assert.NoError(t, err)
		assert.Contains(t, conn, "postgres://u:p@localhost:5432/d")
	})
}
