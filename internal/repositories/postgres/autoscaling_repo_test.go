//go:build integration

package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAutoScalingRepo_Integration(t *testing.T) {
	db := setupDB(t)
	defer db.Close()
	repo := NewAutoScalingRepo(db)
	ctx := context.Background()

	// Cleanup
	_, _ = db.Exec(ctx, "DELETE FROM scaling_group_instances")
	_, _ = db.Exec(ctx, "DELETE FROM scaling_policies")
	_, _ = db.Exec(ctx, "DELETE FROM scaling_groups")
	_, _ = db.Exec(ctx, "DELETE FROM vpcs")
	_, _ = db.Exec(ctx, "DELETE FROM instances")

	vpcID := uuid.New()
	_, err := db.Exec(ctx, "INSERT INTO vpcs (id, name, network_id, created_at) VALUES ($1, $2, $3, $4)",
		vpcID, "asg-vpc", "net-asg", time.Now())
	require.NoError(t, err)

	groupID := uuid.New()
	group := &domain.ScalingGroup{
		ID:             groupID,
		Name:           "test-asg",
		VpcID:          vpcID,
		Image:          "nginx",
		MinInstances:   1,
		MaxInstances:   5,
		DesiredCount:   2,
		Status:         "ACTIVE",
		IdempotencyKey: "asg-key-1",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	t.Run("Scaling Group CRUD", func(t *testing.T) {
		err := repo.CreateGroup(ctx, group)
		require.NoError(t, err)

		fetched, err := repo.GetGroupByID(ctx, groupID)
		require.NoError(t, err)
		assert.Equal(t, group.Name, fetched.Name)

		fetched2, err := repo.GetGroupByIdempotencyKey(ctx, "asg-key-1")
		require.NoError(t, err)
		assert.Equal(t, groupID, fetched2.ID)

		group.DesiredCount = 3
		err = repo.UpdateGroup(ctx, group)
		require.NoError(t, err)

		fetched3, err := repo.GetGroupByID(ctx, groupID)
		require.NoError(t, err)
		assert.Equal(t, 3, fetched3.DesiredCount)

		count, err := repo.CountGroupsByVPC(ctx, vpcID)
		require.NoError(t, err)
		assert.Equal(t, 1, count)

		list, err := repo.ListGroups(ctx)
		require.NoError(t, err)
		assert.NotEmpty(t, list)
	})

	t.Run("Policy Management", func(t *testing.T) {
		policyID := uuid.New()
		policy := &domain.ScalingPolicy{
			ID:             policyID,
			ScalingGroupID: groupID,
			Name:           "scale-out",
			MetricType:     "cpu",
			TargetValue:    70.0,
			ScaleOutStep:   1,
			ScaleInStep:    1,
			CooldownSec:    300,
		}

		err := repo.CreatePolicy(ctx, policy)
		require.NoError(t, err)

		policies, err := repo.GetPoliciesForGroup(ctx, groupID)
		require.NoError(t, err)
		assert.Len(t, policies, 1)
		assert.Equal(t, "scale-out", policies[0].Name)

		now := time.Now()
		err = repo.UpdatePolicyLastScaled(ctx, policyID, now)
		require.NoError(t, err)

		policyMap, err := repo.GetAllPolicies(ctx, []uuid.UUID{groupID})
		require.NoError(t, err)
		assert.NotNil(t, policyMap[groupID][0].LastScaledAt)

		err = repo.DeletePolicy(ctx, policyID)
		require.NoError(t, err)

		policies, err = repo.GetPoliciesForGroup(ctx, groupID)
		require.NoError(t, err)
		assert.Empty(t, policies)
	})

	t.Run("Group Instance Management", func(t *testing.T) {
		instID := uuid.New()
		_, err := db.Exec(ctx, "INSERT INTO instances (id, name, image, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			instID, "asg-inst", "nginx", "running", 1, time.Now(), time.Now())
		require.NoError(t, err)

		err = repo.AddInstanceToGroup(ctx, groupID, instID)
		require.NoError(t, err)

		ids, err := repo.GetInstancesInGroup(ctx, groupID)
		require.NoError(t, err)
		assert.Contains(t, ids, instID)

		instMap, err := repo.GetAllScalingGroupInstances(ctx, []uuid.UUID{groupID})
		require.NoError(t, err)
		assert.Contains(t, instMap[groupID], instID)

		err = repo.RemoveInstanceFromGroup(ctx, groupID, instID)
		require.NoError(t, err)

		ids, err = repo.GetInstancesInGroup(ctx, groupID)
		require.NoError(t, err)
		assert.Empty(t, ids)
	})

	t.Run("Metrics", func(t *testing.T) {
		instID := uuid.New()
		_, _ = db.Exec(ctx, "INSERT INTO instances (id, name, image, status, version, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			instID, "metric-inst", "nginx", "running", 1, time.Now(), time.Now())

		_, err := db.Exec(ctx, "INSERT INTO metrics_history (instance_id, cpu_percent, memory_bytes, recorded_at) VALUES ($1, $2, $3, $4)",
			instID, 50.0, 1024*1024*100, time.Now().Add(-1*time.Minute))
		require.NoError(t, err)

		avg, err := repo.GetAverageCPU(ctx, []uuid.UUID{instID}, time.Now().Add(-5*time.Minute))
		require.NoError(t, err)
		assert.InDelta(t, 50.0, avg, 0.1)
	})

	t.Run("Delete Group", func(t *testing.T) {
		err := repo.DeleteGroup(ctx, groupID)
		require.NoError(t, err)

		_, err = repo.GetGroupByID(ctx, groupID)
		assert.Error(t, err)
	})
}
