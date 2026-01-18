package k8s

import (
	"context"
	"time"

	"github.com/poyrazk/thecloud/internal/core/domain"
)

// MockProvisioner is a simulation provisioner for testing.
type MockProvisioner struct{}

func NewMockProvisioner() *MockProvisioner {
	return &MockProvisioner{}
}

func (p *MockProvisioner) Provision(ctx context.Context, cluster *domain.Cluster) error {
	// Simulate work
	time.Sleep(2 * time.Second)
	cluster.ControlPlaneIPs = []string{"10.0.0.10"}
	cluster.Kubeconfig = "apiVersion: v1\nclusters:\n- cluster:\n    server: https://10.0.0.10:6443\n  name: mock-cluster"
	return nil
}

func (p *MockProvisioner) Deprovision(ctx context.Context, cluster *domain.Cluster) error {
	// Simulate work
	time.Sleep(1 * time.Second)
	return nil
}

func (p *MockProvisioner) GetStatus(ctx context.Context, cluster *domain.Cluster) (domain.ClusterStatus, error) {
	return domain.ClusterStatusRunning, nil
}
