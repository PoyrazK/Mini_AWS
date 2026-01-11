package ports

import (
	"context"
	"io"
)

// ComputeBackend abstracts the underlying infrastructure provider (Docker, Libvirt, etc.)
type ComputeBackend interface {
	// Instance Lifecycle
	CreateInstance(ctx context.Context, opts CreateInstanceOptions) (string, error)
	StopInstance(ctx context.Context, id string) error
	DeleteInstance(ctx context.Context, id string) error
	GetInstanceLogs(ctx context.Context, id string) (io.ReadCloser, error)
	GetInstanceStats(ctx context.Context, id string) (io.ReadCloser, error)
	GetInstancePort(ctx context.Context, id string, internalPort string) (int, error)
	GetInstanceIP(ctx context.Context, id string) (string, error)
	GetConsoleURL(ctx context.Context, id string) (string, error)

	// Execution
	Exec(ctx context.Context, id string, cmd []string) (string, error)
	RunTask(ctx context.Context, opts RunTaskOptions) (string, error)
	WaitTask(ctx context.Context, id string) (int64, error)

	// Network Management
	CreateNetwork(ctx context.Context, name string) (string, error)
	DeleteNetwork(ctx context.Context, id string) error

	// Volume/Disk Attachment (Physical/Block)
	AttachVolume(ctx context.Context, id string, volumePath string) error
	DetachVolume(ctx context.Context, id string, volumePath string) error

	// Health
	Ping(ctx context.Context) error
	Type() string
}
