package workers

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/core/services"
)

type ProvisionWorker struct {
	instSvc   *services.InstanceService
	taskQueue ports.TaskQueue
	logger    *slog.Logger
}

func NewProvisionWorker(instSvc *services.InstanceService, taskQueue ports.TaskQueue, logger *slog.Logger) *ProvisionWorker {
	return &ProvisionWorker{
		instSvc:   instSvc,
		taskQueue: taskQueue,
		logger:    logger,
	}
}

func (w *ProvisionWorker) Run(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	w.logger.Info("starting provision worker")

	for {
		select {
		case <-ctx.Done():
			w.logger.Info("stopping provision worker")
			return
		default:
			// Dequeue task
			msg, err := w.taskQueue.Dequeue(ctx, "provision_queue")
			if err != nil {
				// redis.Nil or other error
				time.Sleep(1 * time.Second)
				continue
			}

			if msg == "" {
				continue
			}

			var job domain.ProvisionJob
			if err := json.Unmarshal([]byte(msg), &job); err != nil {
				w.logger.Error("failed to unmarshal provision job", "error", err)
				continue
			}

			w.logger.Info("processing provision job", "instance_id", job.InstanceID)

			// Process job
			// We use a new context for the actual provisioning so it's not canceled by the worker loop immediate next cycle
			// though BRPop blocks anyway.
			if err := w.instSvc.Provision(context.Background(), job.InstanceID, job.Volumes); err != nil {
				w.logger.Error("failed to provision instance", "instance_id", job.InstanceID, "error", err)
			} else {
				w.logger.Info("successfully provisioned instance", "instance_id", job.InstanceID)
			}
		}
	}
}
