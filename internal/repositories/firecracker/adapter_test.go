//go:build linux

package firecracker

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFirecrackerAdapter_InterfaceCompliance(t *testing.T) {
	var _ ports.ComputeBackend = (*FirecrackerAdapter)(nil)
}

func TestNewFirecrackerAdapter(t *testing.T) {
	logger := slog.Default()
	cfg := Config{
		BinaryPath: "/usr/local/bin/firecracker",
		KernelPath: "/var/lib/thecloud/vmlinux",
		RootfsPath: "/var/lib/thecloud/rootfs.ext4",
		SocketDir:  "/tmp/firecracker-test",
	}

	adapter, err := NewFirecrackerAdapter(logger, cfg)
	require.NoError(t, err)

	assert.NotNil(t, adapter)
	assert.Equal(t, "firecracker", adapter.Type())
	assert.Equal(t, cfg.SocketDir, adapter.cfg.SocketDir)

	t.Run("InvalidSocketDir", func(t *testing.T) {
		// Use a path that is likely to fail on linux (permission denied or similar)
		// Or a path that cannot be created because a file exists with the same name.
		tmpFile, err := os.CreateTemp("", "fc-not-a-dir")
		require.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = NewFirecrackerAdapter(logger, Config{SocketDir: tmpFile.Name()})
		assert.Error(t, err)
	})
}

func TestFirecrackerAdapter_DeleteInstance(t *testing.T) {
	logger := slog.Default()
	cfg := Config{
		SocketDir: t.TempDir(),
		MockMode:  true,
	}
	adapter, err := NewFirecrackerAdapter(logger, cfg)
	require.NoError(t, err)

	ctx := context.Background()

	t.Run("InvalidID", func(t *testing.T) {
		err := adapter.DeleteInstance(ctx, "../invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid instance ID format")
	})

	t.Run("NonExistentInstance", func(t *testing.T) {
		err := adapter.DeleteInstance(ctx, "nonexistent")
		assert.NoError(t, err) // Should return nil if not found
	})
}

func TestFirecrackerAdapter_WaitTask_Mock(t *testing.T) {
	logger := slog.Default()
	cfg := Config{
		MockMode: true,
	}
	adapter, err := NewFirecrackerAdapter(logger, cfg)
	require.NoError(t, err)

	ctx := context.Background()
	exitCode, err := adapter.WaitTask(ctx, "any")
	assert.NoError(t, err)
	assert.Equal(t, int64(0), exitCode)
}
