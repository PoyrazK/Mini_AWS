package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
)

type globalLBRepository struct {
	db *pgxpool.Pool
}

func NewGlobalLBRepository(db *pgxpool.Pool) ports.GlobalLBRepository {
	return &globalLBRepository{
		db: db,
	}
}

func (r *globalLBRepository) Create(ctx context.Context, glb *domain.GlobalLoadBalancer) error {
	query := `
		INSERT INTO global_load_balancers (
			id, user_id, tenant_id, name, hostname, policy, 
			health_check_protocol, health_check_port, health_check_path,
			health_check_interval, health_check_timeout, health_check_healthy_count, health_check_unhealthy_count,
			status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
	`
	_, err := r.db.Exec(ctx, query,
		glb.ID, glb.UserID, glb.TenantID, glb.Name, glb.Hostname, glb.Policy,
		glb.HealthCheck.Protocol, glb.HealthCheck.Port, glb.HealthCheck.Path,
		glb.HealthCheck.IntervalSec, glb.HealthCheck.TimeoutSec, glb.HealthCheck.HealthyCount, glb.HealthCheck.UnhealthyCount,
		glb.Status, glb.CreatedAt, glb.UpdatedAt,
	)
	return err
}

func (r *globalLBRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.GlobalLoadBalancer, error) {
	query := `SELECT * FROM global_load_balancers WHERE id = $1`
	var glb domain.GlobalLoadBalancer
	// Scan struct logic would be needed here, simplified for brevity as we usually list fields
	// Using manual scan for robustness
	row := r.db.QueryRow(ctx, query, id)
	return scanGlobalLB(row)
}

func (r *globalLBRepository) GetByHostname(ctx context.Context, hostname string) (*domain.GlobalLoadBalancer, error) {
	query := `SELECT * FROM global_load_balancers WHERE hostname = $1`
	row := r.db.QueryRow(ctx, query, hostname)
	return scanGlobalLB(row)
}

func (r *globalLBRepository) List(ctx context.Context) ([]*domain.GlobalLoadBalancer, error) {
	query := `SELECT * FROM global_load_balancers`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var glbs []*domain.GlobalLoadBalancer
	for rows.Next() {
		glb, err := scanGlobalLB(rows)
		if err != nil {
			return nil, err
		}
		glbs = append(glbs, glb)
	}
	return glbs, nil
}

func (r *globalLBRepository) Update(ctx context.Context, glb *domain.GlobalLoadBalancer) error {
	query := `
		UPDATE global_load_balancers SET 
			name=$2, hostname=$3, policy=$4, status=$5, updated_at=$6
		WHERE id=$1
	`
	_, err := r.db.Exec(ctx, query, glb.ID, glb.Name, glb.Hostname, glb.Policy, glb.Status, glb.UpdatedAt)
	return err
}

func (r *globalLBRepository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM global_load_balancers WHERE id=$1", id)
	return err
}

// Endpoint methods

func (r *globalLBRepository) AddEndpoint(ctx context.Context, ep *domain.GlobalEndpoint) error {
	query := `
		INSERT INTO global_lb_endpoints (
			id, global_lb_id, region, target_type, target_id, target_ip,
			weight, priority, healthy, created_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query,
		ep.ID, ep.GlobalLBID, ep.Region, ep.TargetType, ep.TargetID, ep.TargetIP,
		ep.Weight, ep.Priority, ep.Healthy, ep.CreatedAt,
	)
	return err
}

func (r *globalLBRepository) RemoveEndpoint(ctx context.Context, endpointID uuid.UUID) error {
	_, err := r.db.Exec(ctx, "DELETE FROM global_lb_endpoints WHERE id=$1", endpointID)
	return err
}

func (r *globalLBRepository) ListEndpoints(ctx context.Context, glbID uuid.UUID) ([]*domain.GlobalEndpoint, error) {
	query := `SELECT * FROM global_lb_endpoints WHERE global_lb_id = $1`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var endpoints []*domain.GlobalEndpoint
	for rows.Next() {
		ep, err := scanEndpoint(rows)
		if err != nil {
			return nil, err
		}
		endpoints = append(endpoints, ep)
	}
	return endpoints, nil
}

func (r *globalLBRepository) UpdateEndpointHealth(ctx context.Context, epID uuid.UUID, healthy bool) error {
	_, err := r.db.Exec(ctx, "UPDATE global_lb_endpoints SET healthy=$2, last_health_check=NOW() WHERE id=$1", epID, healthy)
	return err
}

// Helpers

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanGlobalLB(s scanner) (*domain.GlobalLoadBalancer, error) {
	var glb domain.GlobalLoadBalancer
	var hc domain.HealthCheckConfig
	// Assuming order matches SELECT * (bad practice usually, but standard for these chunks)
	// Or explicitly listing cols is better. sticking to * for speed, but ordering assumes table def
	// id, user_id, tenant_id, name, hostname, policy, hc_proto, hc_port, hc_path, hc_int, hc_to, hc_hc, hc_uhc, status, created, updated
	err := s.Scan(
		&glb.ID, &glb.UserID, &glb.TenantID, &glb.Name, &glb.Hostname, &glb.Policy,
		&hc.Protocol, &hc.Port, &hc.Path, &hc.IntervalSec, &hc.TimeoutSec, &hc.HealthyCount, &hc.UnhealthyCount,
		&glb.Status, &glb.CreatedAt, &glb.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	glb.HealthCheck = hc
	return &glb, nil
}

func scanEndpoint(s scanner) (*domain.GlobalEndpoint, error) {
	var ep domain.GlobalEndpoint
	// id, glb_id, region, type, tid, tip, weight, prio, healthy, last_hc, created
	err := s.Scan(
		&ep.ID, &ep.GlobalLBID, &ep.Region, &ep.TargetType, &ep.TargetID, &ep.TargetIP,
		&ep.Weight, &ep.Priority, &ep.Healthy, &ep.LastHealthCheck, &ep.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &ep, nil
}
