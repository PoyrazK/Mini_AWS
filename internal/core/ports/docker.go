package ports

import (
	"context"
	"io"
)

// DockerClient defines the interface for interacting with the container engine.
type DockerClient interface {
	CreateContainer(ctx context.Context, name, image string, ports []string) (string, error)
	StopContainer(ctx context.Context, containerID string) error
	RemoveContainer(ctx context.Context, containerID string) error
	GetLogs(ctx context.Context, containerID string) (io.ReadCloser, error)
}
