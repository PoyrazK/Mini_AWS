package services_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestElasticIPService_Unit(t *testing.T) {
	mockRepo := new(MockElasticIPRepo)
	mockInstRepo := new(MockInstanceRepo)
	mockAuditSvc := new(MockAuditService)
	
	params := services.ElasticIPServiceParams{
		Repo:         mockRepo,
		InstanceRepo: mockInstRepo,
		AuditSvc:     mockAuditSvc,
		Logger:       slog.Default(),
	}
	svc := services.NewElasticIPService(params)

	ctx := context.Background()
	userID := uuid.New()
	tenantID := uuid.New()
	ctx = appcontext.WithUserID(ctx, userID)
	ctx = appcontext.WithTenantID(ctx, tenantID)

	t.Run("AllocateIP", func(t *testing.T) {
		mockRepo.On("Create", mock.Anything, mock.Anything).Return(nil).Once()
		mockAuditSvc.On("Log", mock.Anything, userID, "eip.allocate", "eip", mock.Anything, mock.Anything).Return(nil).Once()

		eip, err := svc.AllocateIP(ctx)
		assert.NoError(t, err)
		assert.NotNil(t, eip)
		assert.Equal(t, userID, eip.UserID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("AssociateIP_Success", func(t *testing.T) {
		eipID := uuid.New()
		instID := uuid.New()
		vpcID := uuid.New()
		
		eip := &domain.ElasticIP{ID: eipID, UserID: userID, Status: domain.EIPStatusAllocated}
		inst := &domain.Instance{ID: instID, VpcID: &vpcID, Status: domain.StatusRunning}
		
		mockRepo.On("GetByID", mock.Anything, eipID).Return(eip, nil).Once()
		mockInstRepo.On("GetByID", mock.Anything, instID).Return(inst, nil).Once()
		mockRepo.On("GetByInstanceID", mock.Anything, instID).Return(nil, errors.New(errors.NotFound, "not found")).Once()
		mockRepo.On("Update", mock.Anything, mock.Anything).Return(nil).Once()
		mockAuditSvc.On("Log", mock.Anything, userID, "eip.associate", "eip", eipID.String(), mock.Anything).Return(nil).Once()

		res, err := svc.AssociateIP(ctx, eipID, instID)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, domain.EIPStatusAssociated, res.Status)
	})
}
