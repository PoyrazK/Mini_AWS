package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/errors"
)

// ClusterRepository implements ports.ClusterRepository.
type ClusterRepository struct {
	db DB
}

// NewClusterRepository constructs a new ClusterRepository.
func NewClusterRepository(db DB) *ClusterRepository {
	return &ClusterRepository{db: db}
}

func (r *ClusterRepository) Create(ctx context.Context, cluster *domain.Cluster) error {
	query := `
		INSERT INTO clusters (id, user_id, vpc_id, name, version, control_plane_ips, worker_count, status, kubeconfig, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(ctx, query,
		cluster.ID, cluster.UserID, cluster.VpcID, cluster.Name, cluster.Version, cluster.ControlPlaneIPs, cluster.WorkerCount,
		string(cluster.Status), cluster.Kubeconfig, cluster.CreatedAt, cluster.UpdatedAt,
	)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to create cluster", err)
	}
	return nil
}

func (r *ClusterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Cluster, error) {
	userID := appcontext.UserIDFromContext(ctx)
	query := `
		SELECT id, user_id, vpc_id, name, version, COALESCE(control_plane_ips, '{}'), worker_count, status, COALESCE(kubeconfig, ''), created_at, updated_at
		FROM clusters
		WHERE id = $1 AND user_id = $2
	`
	return r.scanCluster(r.db.QueryRow(ctx, query, id, userID))
}

func (r *ClusterRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Cluster, error) {
	query := `
		SELECT id, user_id, vpc_id, name, version, COALESCE(control_plane_ips, '{}'), worker_count, status, COALESCE(kubeconfig, ''), created_at, updated_at
		FROM clusters
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to list clusters", err)
	}
	defer rows.Close()

	var clusters []*domain.Cluster
	for rows.Next() {
		c, err := r.scanCluster(rows)
		if err != nil {
			return nil, err
		}
		clusters = append(clusters, c)
	}
	return clusters, nil
}

func (r *ClusterRepository) Update(ctx context.Context, cluster *domain.Cluster) error {
	query := `
		UPDATE clusters
		SET vpc_id = $1, name = $2, version = $3, control_plane_ips = $4, worker_count = $5, status = $6, kubeconfig = $7, updated_at = $8
		WHERE id = $9 AND user_id = $10
	`
	_, err := r.db.Exec(ctx, query,
		cluster.VpcID, cluster.Name, cluster.Version, cluster.ControlPlaneIPs, cluster.WorkerCount,
		string(cluster.Status), cluster.Kubeconfig, time.Now(), cluster.ID, cluster.UserID,
	)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to update cluster", err)
	}
	return nil
}

func (r *ClusterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	userID := appcontext.UserIDFromContext(ctx)
	query := `DELETE FROM clusters WHERE id = $1 AND user_id = $2`
	_, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to delete cluster", err)
	}
	return nil
}

func (r *ClusterRepository) AddNode(ctx context.Context, node *domain.ClusterNode) error {
	query := `
		INSERT INTO cluster_nodes (id, cluster_id, instance_id, role, status, joined_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, node.ID, node.ClusterID, node.InstanceID, string(node.Role), node.Status, node.JoinedAt)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to add cluster node", err)
	}
	return nil
}

func (r *ClusterRepository) GetNodes(ctx context.Context, clusterID uuid.UUID) ([]*domain.ClusterNode, error) {
	query := `
		SELECT id, cluster_id, instance_id, role, status, joined_at
		FROM cluster_nodes
		WHERE cluster_id = $1
	`
	rows, err := r.db.Query(ctx, query, clusterID)
	if err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to get cluster nodes", err)
	}
	defer rows.Close()

	var nodes []*domain.ClusterNode
	for rows.Next() {
		var n domain.ClusterNode
		var role string
		if err := rows.Scan(&n.ID, &n.ClusterID, &n.InstanceID, &role, &n.Status, &n.JoinedAt); err != nil {
			return nil, errors.Wrap(errors.Internal, "failed to scan cluster node", err)
		}
		n.Role = domain.NodeRole(role)
		nodes = append(nodes, &n)
	}
	return nodes, nil
}

func (r *ClusterRepository) DeleteNode(ctx context.Context, nodeID uuid.UUID) error {
	query := `DELETE FROM cluster_nodes WHERE id = $1`
	_, err := r.db.Exec(ctx, query, nodeID)
	if err != nil {
		return errors.Wrap(errors.Internal, "failed to delete cluster node", err)
	}
	return nil
}

func (r *ClusterRepository) scanCluster(row pgx.Row) (*domain.Cluster, error) {
	var c domain.Cluster
	var status string
	err := row.Scan(&c.ID, &c.UserID, &c.VpcID, &c.Name, &c.Version, &c.ControlPlaneIPs, &c.WorkerCount, &status, &c.Kubeconfig, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, errors.Wrap(errors.Internal, "failed to scan cluster", err)
	}
	c.Status = domain.ClusterStatus(status)
	return &c, nil
}
