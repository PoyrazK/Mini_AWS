package services

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	appcontext "github.com/poyrazk/thecloud/internal/core/context"
	"github.com/poyrazk/thecloud/internal/core/domain"
	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/poyrazk/thecloud/internal/errors"
)

type SnapshotService struct {
	repo       ports.SnapshotRepository
	volumeRepo ports.VolumeRepository
	compute    ports.ComputeBackend
	eventSvc   ports.EventService
	auditSvc   ports.AuditService
	logger     *slog.Logger
	baseDir    string
}

func NewSnapshotService(
	repo ports.SnapshotRepository,
	volumeRepo ports.VolumeRepository,
	compute ports.ComputeBackend,
	eventSvc ports.EventService,
	auditSvc ports.AuditService,
	logger *slog.Logger,
) *SnapshotService {
	baseDir := "./thecloud-data/snapshots"
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		logger.Error("failed to create snapshots directory", "error", err)
	}

	return &SnapshotService{
		repo:       repo,
		volumeRepo: volumeRepo,
		compute:    compute,
		eventSvc:   eventSvc,
		auditSvc:   auditSvc,
		logger:     logger,
		baseDir:    baseDir,
	}
}

func (s *SnapshotService) CreateSnapshot(ctx context.Context, volumeID uuid.UUID, description string) (*domain.Snapshot, error) {
	// 1. Get volume
	vol, err := s.volumeRepo.GetByID(ctx, volumeID)
	if err != nil {
		return nil, err
	}

	// 2. Create domain entity
	snapshot := &domain.Snapshot{
		ID:          uuid.New(),
		UserID:      appcontext.UserIDFromContext(ctx),
		VolumeID:    volumeID,
		VolumeName:  vol.Name,
		SizeGB:      vol.SizeGB,
		Status:      domain.SnapshotStatusCreating,
		Description: description,
		CreatedAt:   time.Now(),
	}

	// 3. Persist to DB initially
	if err := s.repo.Create(ctx, snapshot); err != nil {
		return nil, err
	}

	// Return a copy to avoid data race with the background goroutine
	snapshotCopy := *snapshot

	// 4. Perform async snapshot (for now we'll do it synchronously or simulate async)
	go func() {
		// Use a fresh context for background task
		bgCtx := context.Background()
		err := s.performSnapshot(bgCtx, vol, snapshot)
		if err != nil {
			s.logger.Error("failed to perform snapshot", "snapshot_id", snapshot.ID, "error", err)
			snapshot.Status = domain.SnapshotStatusError
		} else {
			snapshot.Status = domain.SnapshotStatusAvailable
		}
		_ = s.repo.Update(bgCtx, snapshot)
	}()

	_ = s.eventSvc.RecordEvent(ctx, "SNAPSHOT_CREATE", snapshot.ID.String(), "SNAPSHOT", map[string]interface{}{
		"volume_id": volumeID.String(),
	})

	_ = s.auditSvc.Log(ctx, snapshot.UserID, "snapshot.create", "snapshot", snapshot.ID.String(), map[string]interface{}{
		"volume_id": volumeID.String(),
	})

	return &snapshotCopy, nil
}

func (s *SnapshotService) performSnapshot(ctx context.Context, vol *domain.Volume, snapshot *domain.Snapshot) error {
	dockerVolName := "thecloud-vol-" + vol.ID.String()[:8]
	snapshotPath, _ := filepath.Abs(filepath.Join(s.baseDir, snapshot.ID.String()+".tar.gz"))

	// Ensure parent dir exists (it should, but just in case)
	_ = os.MkdirAll(filepath.Dir(snapshotPath), 0755)

	if err := s.compute.CreateVolumeSnapshot(ctx, dockerVolName, snapshotPath); err != nil {
		return fmt.Errorf("failed to create volume snapshot: %w", err)
	}

	return nil
}

func (s *SnapshotService) ListSnapshots(ctx context.Context) ([]*domain.Snapshot, error) {
	userID := appcontext.UserIDFromContext(ctx)
	return s.repo.ListByUserID(ctx, userID)
}

func (s *SnapshotService) GetSnapshot(ctx context.Context, id uuid.UUID) (*domain.Snapshot, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *SnapshotService) DeleteSnapshot(ctx context.Context, id uuid.UUID) error {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 1. Delete file
	snapshotPath := filepath.Join(s.baseDir, snapshot.ID.String()+".tar.gz")
	_ = os.Remove(snapshotPath)

	// 2. Delete from DB
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	_ = s.eventSvc.RecordEvent(ctx, "SNAPSHOT_DELETE", id.String(), "SNAPSHOT", map[string]interface{}{})
	_ = s.auditSvc.Log(ctx, snapshot.UserID, "snapshot.delete", "snapshot", id.String(), map[string]interface{}{})

	return nil
}

func (s *SnapshotService) RestoreSnapshot(ctx context.Context, snapshotID uuid.UUID, newVolumeName string) (*domain.Volume, error) {
	snapshot, err := s.repo.GetByID(ctx, snapshotID)
	if err != nil {
		return nil, err
	}

	if snapshot.Status != domain.SnapshotStatusAvailable {
		return nil, errors.New(errors.InvalidInput, "cannot restore from snapshot that is not available")
	}

	// 1. Create new volume domain entity
	vol := &domain.Volume{
		ID:        uuid.New(),
		UserID:    snapshot.UserID,
		Name:      newVolumeName,
		SizeGB:    snapshot.SizeGB,
		Status:    domain.VolumeStatusAvailable,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 2. Create Docker Volume
	dockerVolName := "thecloud-vol-" + vol.ID.String()[:8]
	if err := s.compute.CreateVolume(ctx, dockerVolName); err != nil {
		return nil, errors.Wrap(errors.Internal, "failed to create docker volume for restore", err)
	}

	// 3. Extract snapshot into new volume
	snapshotPath, _ := filepath.Abs(filepath.Join(s.baseDir, snapshot.ID.String()+".tar.gz"))

	if err := s.compute.RestoreVolumeSnapshot(ctx, dockerVolName, snapshotPath); err != nil {
		_ = s.compute.DeleteVolume(ctx, dockerVolName)
		return nil, fmt.Errorf("failed to restore volume snapshot: %w", err)
	}

	// 4. Save to DB
	if err := s.volumeRepo.Create(ctx, vol); err != nil {
		_ = s.compute.DeleteVolume(ctx, dockerVolName)
		return nil, err
	}

	_ = s.eventSvc.RecordEvent(ctx, "VOLUME_RESTORE", vol.ID.String(), "VOLUME", map[string]interface{}{
		"snapshot_id": snapshotID.String(),
	})

	_ = s.auditSvc.Log(ctx, vol.UserID, "volume.restore", "volume", vol.ID.String(), map[string]interface{}{
		"snapshot_id": snapshotID.String(),
	})

	return vol, nil
}
