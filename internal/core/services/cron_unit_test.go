package services_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCronService_Unit(t *testing.T) {
	repo := new(MockCronRepository)
	eventSvc := new(MockEventService)
	auditSvc := new(MockAuditService)
	svc := services.NewCronService(repo, eventSvc, auditSvc)
	
	ctx := context.Background()
	userID := uuid.New()
	ctx = appcontext.WithUserID(ctx, userID)

	t.Run("CreateJob", func(t *testing.T) {
		repo.On("CreateJob", mock.Anything, mock.Anything).Return(nil).Once()
		eventSvc.On("RecordEvent", mock.Anything, "CRON_JOB_CREATED", mock.Anything, "CRON_JOB", mock.Anything).Return(nil).Once()
		auditSvc.On("Log", mock.Anything, userID, "cron.job_create", "cron_job", mock.Anything, mock.Anything).Return(nil).Once()

		job, err := svc.CreateJob(ctx, "test-job", "*/5 * * * *", "http://test.com", "POST", "{}")
		assert.NoError(t, err)
		assert.NotNil(t, job)
		repo.AssertExpectations(t)
	})

	t.Run("PauseResume", func(t *testing.T) {
		jobID := uuid.New()
		job := &domain.CronJob{ID: jobID, UserID: userID, Schedule: "*/5 * * * *"}
		
		repo.On("GetJobByID", mock.Anything, jobID, userID).Return(job, nil).Twice()
		repo.On("UpdateJob", mock.Anything, mock.MatchedBy(func(j *domain.CronJob) bool {
			return j.Status == domain.CronStatusPaused
		})).Return(nil).Once()
		repo.On("UpdateJob", mock.Anything, mock.MatchedBy(func(j *domain.CronJob) bool {
			return j.Status == domain.CronStatusActive
		})).Return(nil).Once()

		err := svc.PauseJob(ctx, jobID)
		assert.NoError(t, err)
		err = svc.ResumeJob(ctx, jobID)
		assert.NoError(t, err)
	})

	t.Run("DeleteJob", func(t *testing.T) {
		jobID := uuid.New()
		repo.On("GetJobByID", mock.Anything, jobID, userID).Return(&domain.CronJob{ID: jobID}, nil).Once()
		repo.On("DeleteJob", mock.Anything, jobID).Return(nil).Once()
		auditSvc.On("Log", mock.Anything, userID, "cron.job_delete", "cron_job", jobID.String(), mock.Anything).Return(nil).Once()

		err := svc.DeleteJob(ctx, jobID)
		assert.NoError(t, err)
	})
}
