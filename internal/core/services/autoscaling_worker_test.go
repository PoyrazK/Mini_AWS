package services_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyraz/cloud/internal/core/domain"
	"github.com/poyraz/cloud/internal/core/services"
	"github.com/stretchr/testify/mock"
)

func TestAutoScalingWorker_Logic(t *testing.T) {
	ctx := context.Background()
	groupID := uuid.New()
	vpcID := uuid.New()
	lbID := uuid.New()
	now := time.Now()

	t.Run("Scale Out when current < desired", func(t *testing.T) {
		asgRepo, instSvc, lbSvc, eventSvc, clock := newMockWorkerDeps()
		worker := services.NewAutoScalingWorker(asgRepo, instSvc, lbSvc, eventSvc, clock)

		group := &domain.ScalingGroup{
			ID:             groupID,
			Name:           "test-asg",
			VpcID:          vpcID,
			LoadBalancerID: &lbID,
			Image:          "nginx",
			Ports:          "80:80",
			MinInstances:   1,
			MaxInstances:   5,
			DesiredCount:   2,
			CurrentCount:   1,
		}

		instances := []uuid.UUID{uuid.New()}

		asgRepo.On("ListGroups", ctx).Return([]*domain.ScalingGroup{group}, nil).Once()
		asgRepo.On("GetAllScalingGroupInstances", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]uuid.UUID{groupID: instances}, nil).Once()
		asgRepo.On("GetAllPolicies", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]*domain.ScalingPolicy{groupID: {}}, nil).Once()

		clock.On("Now").Return(now).Maybe()

		newInstID := uuid.New()
		instSvc.On("LaunchInstance", ctx, mock.Anything, "nginx", "0:80", &vpcID, []domain.VolumeAttachment(nil)).Return(&domain.Instance{ID: newInstID}, nil).Once()
		asgRepo.On("AddInstanceToGroup", ctx, groupID, newInstID).Return(nil).Once()
		lbSvc.On("AddTarget", ctx, lbID, newInstID, 80, 1).Return(nil).Once()
		eventSvc.On("RecordEvent", ctx, "AUTOSCALING_SCALE_OUT", groupID.String(), "SCALING_GROUP", mock.Anything).Return(nil).Once()
		asgRepo.On("UpdateGroup", ctx, mock.Anything).Return(nil).Maybe() // reset failures

		worker.Evaluate(ctx)

		asgRepo.AssertExpectations(t)
		instSvc.AssertExpectations(t)
		lbSvc.AssertExpectations(t)
	})

	t.Run("Scale In when current > desired", func(t *testing.T) {
		asgRepo, instSvc, lbSvc, eventSvc, clock := newMockWorkerDeps()
		worker := services.NewAutoScalingWorker(asgRepo, instSvc, lbSvc, eventSvc, clock)

		instID1 := uuid.New()
		instID2 := uuid.New()
		group := &domain.ScalingGroup{
			ID:             groupID,
			Name:           "test-asg",
			VpcID:          vpcID,
			LoadBalancerID: &lbID,
			MinInstances:   1,
			MaxInstances:   5,
			DesiredCount:   1,
			CurrentCount:   2,
		}

		instances := []uuid.UUID{instID1, instID2}

		asgRepo.On("ListGroups", ctx).Return([]*domain.ScalingGroup{group}, nil).Once()
		asgRepo.On("GetAllScalingGroupInstances", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]uuid.UUID{groupID: instances}, nil).Once()
		asgRepo.On("GetAllPolicies", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]*domain.ScalingPolicy{groupID: {}}, nil).Once()

		lbSvc.On("RemoveTarget", ctx, lbID, instID2).Return(nil).Once()
		asgRepo.On("RemoveInstanceFromGroup", ctx, groupID, instID2).Return(nil).Once()
		instSvc.On("TerminateInstance", ctx, instID2.String()).Return(nil).Once()
		eventSvc.On("RecordEvent", ctx, "AUTOSCALING_SCALE_IN", groupID.String(), "SCALING_GROUP", mock.Anything).Return(nil).Once()

		worker.Evaluate(ctx)

		lbSvc.AssertExpectations(t)
		asgRepo.AssertExpectations(t)
		instSvc.AssertExpectations(t)
	})

	t.Run("Policy trigger scale out", func(t *testing.T) {
		asgRepo, instSvc, lbSvc, eventSvc, clock := newMockWorkerDeps()
		worker := services.NewAutoScalingWorker(asgRepo, instSvc, lbSvc, eventSvc, clock)

		group := &domain.ScalingGroup{
			ID:           groupID,
			CurrentCount: 1,
			DesiredCount: 1, // Set to avoid reconciliation before policy check
			MinInstances: 1,
			MaxInstances: 5,
		}
		instanceIDs := []uuid.UUID{uuid.New()}
		policy := &domain.ScalingPolicy{
			ID:           uuid.New(),
			Name:         "cpu-high",
			MetricType:   "cpu",
			TargetValue:  70.0,
			ScaleOutStep: 1,
			CooldownSec:  300,
		}

		asgRepo.On("ListGroups", ctx).Return([]*domain.ScalingGroup{group}, nil).Once()
		asgRepo.On("GetAllScalingGroupInstances", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]uuid.UUID{groupID: instanceIDs}, nil).Once()
		asgRepo.On("GetAllPolicies", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]*domain.ScalingPolicy{groupID: {policy}}, nil).Once()

		clock.On("Now").Return(now).Maybe()
		asgRepo.On("GetAverageCPU", ctx, instanceIDs, mock.Anything).Return(80.0, nil).Once()

		asgRepo.On("UpdateGroup", ctx, mock.MatchedBy(func(g *domain.ScalingGroup) bool {
			return g.DesiredCount == 2
		})).Return(nil).Once()
		asgRepo.On("UpdatePolicyLastScaled", ctx, policy.ID, mock.Anything).Return(nil).Once()

		worker.Evaluate(ctx)

		asgRepo.AssertExpectations(t)
	})

	t.Run("Policy trigger scale in", func(t *testing.T) {
		asgRepo, instSvc, lbSvc, eventSvc, clock := newMockWorkerDeps()
		worker := services.NewAutoScalingWorker(asgRepo, instSvc, lbSvc, eventSvc, clock)

		group := &domain.ScalingGroup{
			ID:           groupID,
			CurrentCount: 2,
			DesiredCount: 2,
			MinInstances: 1,
			MaxInstances: 5,
		}
		instanceIDs := []uuid.UUID{uuid.New(), uuid.New()}
		policy := &domain.ScalingPolicy{
			ID:          uuid.New(),
			MetricType:  "cpu",
			TargetValue: 70.0,
			ScaleInStep: 1,
			CooldownSec: 300,
		}

		asgRepo.On("ListGroups", ctx).Return([]*domain.ScalingGroup{group}, nil).Once()
		asgRepo.On("GetAllScalingGroupInstances", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]uuid.UUID{groupID: instanceIDs}, nil).Once()
		asgRepo.On("GetAllPolicies", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]*domain.ScalingPolicy{groupID: {policy}}, nil).Once()

		clock.On("Now").Return(now).Maybe()
		asgRepo.On("GetAverageCPU", ctx, instanceIDs, mock.Anything).Return(40.0, nil).Once()

		asgRepo.On("UpdateGroup", ctx, mock.MatchedBy(func(g *domain.ScalingGroup) bool {
			return g.DesiredCount == 1
		})).Return(nil).Once()
		asgRepo.On("UpdatePolicyLastScaled", ctx, policy.ID, mock.Anything).Return(nil).Once()

		worker.Evaluate(ctx)

		asgRepo.AssertExpectations(t)
	})

	t.Run("Policy skipped due to cooldown", func(t *testing.T) {
		asgRepo, instSvc, lbSvc, eventSvc, clock := newMockWorkerDeps()
		worker := services.NewAutoScalingWorker(asgRepo, instSvc, lbSvc, eventSvc, clock)

		lastScaled := now.Add(-1 * time.Minute)
		group := &domain.ScalingGroup{
			ID:           groupID,
			CurrentCount: 1,
			DesiredCount: 1,
			MinInstances: 1,
			MaxInstances: 5,
		}
		instanceIDs := []uuid.UUID{uuid.New()}
		policy := &domain.ScalingPolicy{
			ID:           uuid.New(),
			MetricType:   "cpu",
			TargetValue:  70.0,
			ScaleOutStep: 1,
			CooldownSec:  300, // 5 min
			LastScaledAt: &lastScaled,
		}

		asgRepo.On("ListGroups", ctx).Return([]*domain.ScalingGroup{group}, nil).Once()
		asgRepo.On("GetAllScalingGroupInstances", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]uuid.UUID{groupID: instanceIDs}, nil).Once()
		asgRepo.On("GetAllPolicies", ctx, []uuid.UUID{groupID}).Return(map[uuid.UUID][]*domain.ScalingPolicy{groupID: {policy}}, nil).Once()

		clock.On("Now").Return(now).Maybe()
		// No GetAverageCPU should be called because cooldown skips it
		// Actually, current implementation in evaluatePolicies:
		// for _, policy := range policies {
		//    if cooldown... continue
		// }
		// Wait, GetAverageCPU is called BEFORE the loop.
		// So it WILL be called, but the policy evaluation loop will skip the action.
		asgRepo.On("GetAverageCPU", ctx, instanceIDs, mock.Anything).Return(80.0, nil).Once()

		worker.Evaluate(ctx)

		// Assert UpdateGroup NOT called
		asgRepo.AssertNotCalled(t, "UpdateGroup", mock.Anything, mock.Anything)
		asgRepo.AssertExpectations(t)
	})
}

func newMockWorkerDeps() (*MockAutoScalingRepo, *MockInstanceService, *MockLBService, *MockEventService, *MockClock) {
	return new(MockAutoScalingRepo), new(MockInstanceService), new(MockLBService), new(MockEventService), new(MockClock)
}
