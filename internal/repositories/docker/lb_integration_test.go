//go:build integration

package docker

import (
	"context"
	"testing"

	"github.com/docker/docker/client"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
)

func TestLBProxyAdapter_Integration(t *testing.T) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		t.Skip("Docker not available")
	}
	ctx := context.Background()

	instRepo := new(mockInstanceRepo)
	vpcRepo := new(mockVpcRepo)
	adapter := &LBProxyAdapter{
		cli:          cli,
		instanceRepo: instRepo,
		vpcRepo:      vpcRepo,
	}

	lb := &domain.LoadBalancer{
		ID:        uuid.New(),
		Port:      9999, // Use a high port unlikely to be used
		VpcID:     uuid.New(),
		Algorithm: "round-robin",
	}

	// Mock VPC to use bridge network for simplicity in test
	vpcRepo.On("GetByID", ctx, lb.VpcID).Return(&domain.VPC{NetworkID: "bridge"}, nil)

	t.Run("Deploy and Remove Proxy", func(t *testing.T) {
		// Clean up any leftovers
		_ = adapter.RemoveProxy(ctx, lb.ID)

		// Deploy
		containerID, err := adapter.DeployProxy(ctx, lb, nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, containerID)

		// Verify container exists
		_, err = cli.ContainerInspect(ctx, containerID)
		assert.NoError(t, err)

		// Update (reloads nginx)
		err = adapter.UpdateProxyConfig(ctx, lb, nil)
		assert.NoError(t, err)

		// Remove
		err = adapter.RemoveProxy(ctx, lb.ID)
		assert.NoError(t, err)

		// Verify gone
		_, err = cli.ContainerInspect(ctx, containerID)
		assert.Error(t, err)
	})
}
