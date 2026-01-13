package ovs

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/poyrazk/thecloud/internal/core/ports"
	apperrors "github.com/poyrazk/thecloud/internal/errors"
	"github.com/stretchr/testify/require"
)

type fakeCmd struct {
	runErr  error
	out     []byte
	outErr  error
	runHits int
}

func (c *fakeCmd) Run() error {
	c.runHits++
	return c.runErr
}

func (c *fakeCmd) Output() ([]byte, error) {
	if c.outErr != nil {
		return nil, c.outErr
	}
	return c.out, nil
}

type fakeExecer struct {
	lookPath map[string]string
	lookErr  error
	cmd      *fakeCmd
}

func (e *fakeExecer) LookPath(file string) (string, error) {
	if e.lookErr != nil {
		return "", e.lookErr
	}
	if p, ok := e.lookPath[file]; ok {
		return p, nil
	}
	return "", errors.New("not found")
}

func (e *fakeExecer) CommandContext(ctx context.Context, name string, args ...string) cmd {
	return e.cmd
}

func TestOvsAdapter_CommandErrorsAreWrapped(t *testing.T) {
	fx := &fakeExecer{
		lookPath: map[string]string{"ovs-vsctl": "/bin/ovs-vsctl", "ovs-ofctl": "/bin/ovs-ofctl"},
		cmd:      &fakeCmd{runErr: errors.New("boom")},
	}

	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

	err := a.CreateBridge(context.Background(), "br0", 0)
	require.Error(t, err)
	require.True(t, apperrors.Is(err, apperrors.Internal))
}

func TestOvsAdapter_ListBridges_EmptyOutput(t *testing.T) {
	fx := &fakeExecer{
		cmd: &fakeCmd{out: []byte("\n")},
	}

	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}
	bridges, err := a.ListBridges(context.Background())
	require.NoError(t, err)
	require.Len(t, bridges, 0)
}

func TestOvsAdapter_AddFlowRule_InvalidBridge(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

	err := a.AddFlowRule(context.Background(), "bad bridge", ports.FlowRule{Priority: 1, Match: "ip", Actions: "drop"})
	require.Error(t, err)
	require.True(t, apperrors.Is(err, apperrors.InvalidInput))
}
