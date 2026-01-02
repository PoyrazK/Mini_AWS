package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
)

func TestCreateGroup_SecurityLimits(t *testing.T) {
	mockRepo := new(MockAutoScalingRepo)
	mockVpcRepo := new(MockVpcRepo)
	svc := services.NewAutoScalingService(mockRepo, mockVpcRepo)
	ctx := context.Background()
	vpcID := uuid.New()

	mockVpcRepo.On("GetByID", ctx, vpcID).Return(&domain.VPC{ID: vpcID}, nil)

	t.Run("ExceedsMaxInstances", func(t *testing.T) {
		_, err := svc.CreateGroup(ctx, "test", vpcID, "img", "80:80", 1, 1000, 1, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "max_instances cannot exceed")
	})

	t.Run("ExceedsVPCLimit", func(t *testing.T) {
		mockRepo.On("CountGroupsByVPC", ctx, vpcID).Return(10, nil)
		_, err := svc.CreateGroup(ctx, "test", vpcID, "img", "80:80", 1, 5, 1, nil, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "VPC already has")
	})
}
