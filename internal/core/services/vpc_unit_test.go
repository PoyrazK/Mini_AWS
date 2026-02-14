package services_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestVpcService_Unit(t *testing.T) {
	repo := new(MockVpcRepo)
	lbRepo := new(MockLBRepo)
	network := new(MockNetworkBackend)
	audit := new(MockAuditService)
	logger := slog.Default()

	svc := services.NewVpcService(repo, lbRepo, network, audit, logger, "10.0.0.0/16")

	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.New()
	ctx = appcontext.WithUserID(ctx, userID)
	ctx = appcontext.WithTenantID(ctx, tenantID)

	t.Run("CreateVPC_Success", func(t *testing.T) {
		name := "my-vpc"
		cidr := "10.1.0.0/16"

		network.On("CreateBridge", mock.Anything, mock.Anything, mock.Anything).Return(nil).Once()
		repo.On("Create", mock.Anything, mock.MatchedBy(func(v *domain.VPC) bool {
			return v.Name == name && v.CIDRBlock == cidr && v.UserID == userID
		})).Return(nil).Once()
		audit.On("Log", mock.Anything, userID, "vpc.create", "vpc", mock.Anything, mock.Anything).Return(nil).Once()

		vpc, err := svc.CreateVPC(ctx, name, cidr)
		assert.NoError(t, err)
		assert.NotNil(t, vpc)
		assert.Equal(t, name, vpc.Name)
		repo.AssertExpectations(t)
	})

	t.Run("DeleteVPC_Success", func(t *testing.T) {
		vpcID := uuid.New()
		vpc := &domain.VPC{ID: vpcID, Name: "test-vpc", NetworkID: "br-1", UserID: userID}

		repo.On("GetByID", mock.Anything, vpcID).Return(vpc, nil).Once()
		lbRepo.On("ListAll", mock.Anything).Return([]*domain.LoadBalancer{}, nil).Once()
		network.On("DeleteBridge", mock.Anything, "br-1").Return(nil).Once()
		repo.On("Delete", mock.Anything, vpcID).Return(nil).Once()
		audit.On("Log", mock.Anything, userID, "vpc.delete", "vpc", vpcID.String(), mock.Anything).Return(nil).Once()

		err := svc.DeleteVPC(ctx, vpcID.String())
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
}
