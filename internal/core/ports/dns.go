package ports

import (
	"context"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
)

// DNSRepository manages persistent storage for DNS zones and records.
type DNSRepository interface {
	// Zone operations
	CreateZone(ctx context.Context, zone *domain.DNSZone) error
	GetZoneByID(ctx context.Context, id uuid.UUID) (*domain.DNSZone, error)
	GetZoneByName(ctx context.Context, name string) (*domain.DNSZone, error)
	GetZoneByVPC(ctx context.Context, vpcID uuid.UUID) (*domain.DNSZone, error)
	ListZones(ctx context.Context) ([]*domain.DNSZone, error)
	UpdateZone(ctx context.Context, zone *domain.DNSZone) error
	DeleteZone(ctx context.Context, id uuid.UUID) error

	// Record operations
	CreateRecord(ctx context.Context, record *domain.DNSRecord) error
	GetRecordByID(ctx context.Context, id uuid.UUID) (*domain.DNSRecord, error)
	ListRecordsByZone(ctx context.Context, zoneID uuid.UUID) ([]*domain.DNSRecord, error)
	GetRecordsByInstance(ctx context.Context, instanceID uuid.UUID) ([]*domain.DNSRecord, error)
	UpdateRecord(ctx context.Context, record *domain.DNSRecord) error
	DeleteRecord(ctx context.Context, id uuid.UUID) error
	DeleteRecordsByInstance(ctx context.Context, instanceID uuid.UUID) error
}

// DNSBackend abstracts the actual DNS server (PowerDNS).
type DNSBackend interface {
	// Zone operations in PowerDNS
	CreateZone(ctx context.Context, zoneName string, nameservers []string) error
	DeleteZone(ctx context.Context, zoneName string) error
	GetZone(ctx context.Context, zoneName string) (*ZoneInfo, error)

	// Record operations in PowerDNS
	AddRecords(ctx context.Context, zoneName string, records []RecordSet) error
	UpdateRecords(ctx context.Context, zoneName string, records []RecordSet) error
	DeleteRecords(ctx context.Context, zoneName string, name string, recordType string) error
	ListRecords(ctx context.Context, zoneName string) ([]RecordSet, error)
}

// ZoneInfo represents zone information from PowerDNS.
type ZoneInfo struct {
	Name           string
	Kind           string // Native, Master, Slave
	Serial         uint32
	NotifiedSerial uint32
}

// RecordSet represents a set of records for a name/type in PowerDNS.
type RecordSet struct {
	Name     string // FQDN, e.g., "www.myapp.internal."
	Type     string // A, AAAA, CNAME, etc.
	TTL      int
	Records  []string // Values
	Priority *int     // For MX, SRV
}

// DNSService provides business logic for managing DNS zones and records.
type DNSService interface {
	// Zone operations
	CreateZone(ctx context.Context, vpcID uuid.UUID, name, description string) (*domain.DNSZone, error)
	GetZone(ctx context.Context, idOrName string) (*domain.DNSZone, error)
	GetZoneByVPC(ctx context.Context, vpcID uuid.UUID) (*domain.DNSZone, error)
	ListZones(ctx context.Context) ([]*domain.DNSZone, error)
	DeleteZone(ctx context.Context, idOrName string) error

	// Record operations
	CreateRecord(ctx context.Context, zoneID uuid.UUID, name string, recordType domain.RecordType, content string, ttl int, priority *int) (*domain.DNSRecord, error)
	GetRecord(ctx context.Context, id uuid.UUID) (*domain.DNSRecord, error)
	ListRecords(ctx context.Context, zoneID uuid.UUID) ([]*domain.DNSRecord, error)
	UpdateRecord(ctx context.Context, id uuid.UUID, content string, ttl int, priority *int) (*domain.DNSRecord, error)
	DeleteRecord(ctx context.Context, id uuid.UUID) error

	// Instance auto-registration (called by InstanceService)
	RegisterInstance(ctx context.Context, instance *domain.Instance, ipAddress string) error
	UnregisterInstance(ctx context.Context, instanceID uuid.UUID) error
}
