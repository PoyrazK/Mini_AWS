package ports

import (
	"context"
)

// StorageBackend abstracts block storage operations (LVM, Ceph, etc.)
type StorageBackend interface {
	CreateVolume(ctx context.Context, name string, sizeGB int) (string, error)
	DeleteVolume(ctx context.Context, name string) error
	AttachVolume(ctx context.Context, volumeName, instanceID string) error
	DetachVolume(ctx context.Context, volumeName, instanceID string) error
	CreateSnapshot(ctx context.Context, volumeName, snapshotName string) error
	DeleteSnapshot(ctx context.Context, snapshotName string) error
	Ping(ctx context.Context) error
	Type() string
}
