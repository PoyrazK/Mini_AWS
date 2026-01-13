package workers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/services"
	"github.com/poyrazk/thecloud/internal/repositories/noop"
)

type fakeTaskQueue struct {
	messages []string
	index    int
}

func (f *fakeTaskQueue) Enqueue(ctx context.Context, queueName string, payload interface{}) error {
	return nil
}

func (f *fakeTaskQueue) Dequeue(ctx context.Context, queueName string) (string, error) {
	if f.index < len(f.messages) {
		msg := f.messages[f.index]
		f.index++
		return msg, nil
	}
	return "", nil
}

func TestProvisionWorker_Run(t *testing.T) {
	job := domain.ProvisionJob{
		InstanceID: uuid.New(),
		UserID:     uuid.New(),
		Volumes:    []domain.VolumeAttachment{},
	}
	msg, _ := json.Marshal(job)

	fakeQueue := &fakeTaskQueue{messages: []string{string(msg)}}

	// Create InstanceService with noop dependencies
	instSvc := services.NewInstanceService(services.InstanceServiceParams{
		Repo:       &noop.NoopInstanceRepository{},
		VpcRepo:    &noop.NoopVpcRepository{},
		SubnetRepo: &noop.NoopSubnetRepository{},
		VolumeRepo: &noop.NoopVolumeRepository{},
		Compute:    &noop.NoopComputeBackend{},
		Network:    noop.NewNoopNetworkAdapter(slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))),
		EventSvc:   &noop.NoopEventService{},
		AuditSvc:   &noop.NoopAuditService{},
		TaskQueue:  nil,
		Logger:     slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil)),
	})

	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{Level: slog.LevelDebug})
	logger := slog.New(handler)

	worker := NewProvisionWorker(instSvc, fakeQueue, logger)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)

	go worker.Run(ctx, &wg)

	// Wait a bit for the goroutine to process
	time.Sleep(100 * time.Millisecond)
	cancel()
	wg.Wait()

	// Check that the success log was written
	logOutput := buf.String()
	if !bytes.Contains([]byte(logOutput), []byte("successfully provisioned instance")) {
		t.Errorf("expected success log, got: %s", logOutput)
	}
}
