//go:build !linux

package firecracker

import (
	"log/slog"
	"testing"

	"github.com/poyrazk/thecloud/internal/core/ports"
	"github.com/stretchr/testify/assert"
)

func TestFirecrackerAdapter_InterfaceCompliance(t *testing.T) {
	var _ ports.ComputeBackend = (*FirecrackerAdapter)(nil)
}

func TestNewFirecrackerAdapter(t *testing.T) {
	logger := slog.Default()
	cfg := Config{}

	adapter := NewFirecrackerAdapter(logger, cfg)

	assert.NotNil(t, adapter)
	assert.Equal(t, "firecracker-noop", adapter.Type())
}
