// Package postgres provides PostgreSQL-backed repository implementations.
package postgres

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/errors"
)

// InstanceTypeRepository provides a PostgreSQL implementation for instance types.
type InstanceTypeRepository struct {
	db DB
}

// NewInstanceTypeRepository creates a new InstanceTypeRepository.
func NewInstanceTypeRepository(db DB) *InstanceTypeRepository {
	return &InstanceTypeRepository{db: db}
}

// List returns all available instance types.
func (r *InstanceTypeRepository) List(ctx context.Context) ([]*domain.InstanceType, error) {
	query := `
		SELECT id, name, vcpus, memory_mb, disk_gb, network_mbps, price_per_hour, category
		FROM instance_types
		ORDER BY vcpus ASC, memory_mb ASC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to list instance types", err)
	}
	defer rows.Close()

	var types []*domain.InstanceType
	for rows.Next() {
		t, err := r.scanInstanceType(rows)
		if err != nil {
			return nil, err
		}
		types = append(types, t)
	}
	return types, nil
}

// GetByID retrieves an instance type by its ID.
func (r *InstanceTypeRepository) GetByID(ctx context.Context, id string) (*domain.InstanceType, error) {
	query := `
		SELECT id, name, vcpus, memory_mb, disk_gb, network_mbps, price_per_hour, category
		FROM instance_types
		WHERE id = $1
	`
	return r.scanInstanceType(r.db.QueryRow(ctx, query, id))
}

func (r *InstanceTypeRepository) scanInstanceType(row pgx.Row) (*domain.InstanceType, error) {
	var t domain.InstanceType
	err := row.Scan(
		&t.ID, &t.Name, &t.VCPUs, &t.MemoryMB, &t.DiskGB, &t.NetworkMbps, &t.PricePerHr, &t.Category,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.New(errors.NotFound, "instance type not found")
		}
		return nil, errors.Wrap(errors.Internal, "failed to scan instance type", err)
	}
	return &t, nil
}
