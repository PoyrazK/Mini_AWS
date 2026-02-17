package services_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDNSRepository struct {
	mock.Mock
}

func (m *MockDNSRepository) CreateZone(ctx context.Context, zone *domain.DNSZone) error {
	return m.Called(ctx, zone).Error(0)
}
func (m *MockDNSRepository) GetZoneByID(ctx context.Context, id uuid.UUID) (*domain.DNSZone, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DNSZone), args.Error(1)
}
func (m *MockDNSRepository) GetZoneByName(ctx context.Context, name string) (*domain.DNSZone, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DNSZone), args.Error(1)
}
func (m *MockDNSRepository) GetZoneByVPC(ctx context.Context, vpcID uuid.UUID) (*domain.DNSZone, error) {
	args := m.Called(ctx, vpcID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DNSZone), args.Error(1)
}
func (m *MockDNSRepository) ListZones(ctx context.Context) ([]*domain.DNSZone, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DNSZone), args.Error(1)
}
func (m *MockDNSRepository) UpdateZone(ctx context.Context, zone *domain.DNSZone) error {
	return m.Called(ctx, zone).Error(0)
}
func (m *MockDNSRepository) DeleteZone(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockDNSRepository) CreateRecord(ctx context.Context, record *domain.DNSRecord) error {
	return m.Called(ctx, record).Error(0)
}
func (m *MockDNSRepository) GetRecordByID(ctx context.Context, id uuid.UUID) (*domain.DNSRecord, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.DNSRecord), args.Error(1)
}
func (m *MockDNSRepository) ListRecordsByZone(ctx context.Context, zoneID uuid.UUID) ([]*domain.DNSRecord, error) {
	args := m.Called(ctx, zoneID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DNSRecord), args.Error(1)
}
func (m *MockDNSRepository) GetRecordsByInstance(ctx context.Context, instanceID uuid.UUID) ([]*domain.DNSRecord, error) {
	args := m.Called(ctx, instanceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.DNSRecord), args.Error(1)
}
func (m *MockDNSRepository) UpdateRecord(ctx context.Context, record *domain.DNSRecord) error {
	return m.Called(ctx, record).Error(0)
}
func (m *MockDNSRepository) DeleteRecord(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockDNSRepository) DeleteRecordsByInstance(ctx context.Context, instanceID uuid.UUID) error {
	return m.Called(ctx, instanceID).Error(0)
}

type MockDNSBackend struct {
	mock.Mock
}

func (m *MockDNSBackend) CreateZone(ctx context.Context, name string, nameservers []string) error {
	return m.Called(ctx, name, nameservers).Error(0)
}
func (m *MockDNSBackend) DeleteZone(ctx context.Context, name string) error {
	return m.Called(ctx, name).Error(0)
}
func (m *MockDNSBackend) GetZone(ctx context.Context, name string) (*ports.ZoneInfo, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*ports.ZoneInfo), args.Error(1)
}
func (m *MockDNSBackend) AddRecords(ctx context.Context, zoneID string, records []ports.RecordSet) error {
	return m.Called(ctx, zoneID, records).Error(0)
}
func (m *MockDNSBackend) UpdateRecords(ctx context.Context, zoneID string, records []ports.RecordSet) error {
	return m.Called(ctx, zoneID, records).Error(0)
}
func (m *MockDNSBackend) DeleteRecords(ctx context.Context, zoneID, name, rType string) error {
	return m.Called(ctx, zoneID, name, rType).Error(0)
}
func (m *MockDNSBackend) ListRecords(ctx context.Context, zoneID string) ([]ports.RecordSet, error) {
	args := m.Called(ctx, zoneID)
	return args.Get(0).([]ports.RecordSet), args.Error(1)
}

func TestDNSService_Unit_Extended(t *testing.T) {
	repo := new(MockDNSRepository)
	backend := new(MockDNSBackend)
	vpcRepo := new(MockVpcRepo)
	auditSvc := new(MockAuditService)
	eventSvc := new(MockEventService)
	
	svc := services.NewDNSService(services.DNSServiceParams{
		Repo:     repo,
		Backend:  backend,
		VpcRepo:  vpcRepo,
		AuditSvc: auditSvc,
		EventSvc: eventSvc,
		Logger:   slog.Default(),
	})

	ctx := context.Background()
	userID := uuid.New()
	ctx = appcontext.WithUserID(ctx, userID)

	t.Run("DeleteZone", func(t *testing.T) {
		zoneID := uuid.New()
		zone := &domain.DNSZone{ID: zoneID, Name: "example.com", PowerDNSID: "example.com.", UserID: userID}
		repo.On("GetZoneByID", mock.Anything, zoneID).Return(zone, nil).Once()
		backend.On("DeleteZone", mock.Anything, "example.com.").Return(nil).Once()
		repo.On("DeleteZone", mock.Anything, zoneID).Return(nil).Once()
		auditSvc.On("Log", mock.Anything, userID, "dns.zone.delete", "dns_zone", zoneID.String(), mock.Anything).Return(nil).Once()

		err := svc.DeleteZone(ctx, zoneID.String())
		assert.NoError(t, err)
	})

	t.Run("ListRecords", func(t *testing.T) {
		zoneID := uuid.New()
		repo.On("ListRecordsByZone", mock.Anything, zoneID).Return([]*domain.DNSRecord{}, nil).Once()
		
		res, err := svc.ListRecords(ctx, zoneID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
	})
	
	t.Run("RegisterInstance", func(t *testing.T) {
		inst := &domain.Instance{ID: uuid.New(), Name: "web-1", VpcID: new(uuid.UUID)}
		zone := &domain.DNSZone{ID: uuid.New(), Name: "cluster.local", PowerDNSID: "cluster.local."}
		
		repo.On("GetZoneByVPC", mock.Anything, *inst.VpcID).Return(zone, nil).Once()
		backend.On("AddRecords", mock.Anything, "cluster.local.", mock.Anything).Return(nil).Once()
		repo.On("CreateRecord", mock.Anything, mock.Anything).Return(nil).Once()
		
		err := svc.RegisterInstance(ctx, inst, "10.0.0.1")
		assert.NoError(t, err)
	})
}
