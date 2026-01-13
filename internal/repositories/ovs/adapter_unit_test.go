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

	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}
	bridges, err := a.ListBridges(context.Background())
	require.NoError(t, err)
	require.Len(t, bridges, 0)
}

func TestOvsAdapter_AddFlowRule_InvalidBridge(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

	err := a.AddFlowRule(context.Background(), "bad bridge", ports.FlowRule{Priority: 1, Match: "ip", Actions: "drop"})
	require.Error(t, err)
	require.True(t, apperrors.Is(err, apperrors.InvalidInput))
}

func TestOvsAdapter_AddFlowRule_Success(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

	err := a.AddFlowRule(context.Background(), "br0", ports.FlowRule{Priority: 100, Match: "ip", Actions: "normal"})
	require.NoError(t, err)
}

func TestOvsAdapter_AddPort(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.AddPort(context.Background(), "br0", "port1")
	require.NoError(t, err)
}

func TestOvsAdapter_DeletePort(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.DeletePort(context.Background(), "br0", "port1")
	require.NoError(t, err)
}

func TestOvsAdapter_Ping(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.Ping(context.Background())
	require.NoError(t, err)
}

func TestOvsAdapter_Type(t *testing.T) {
	a := &OvsAdapter{}
	require.Equal(t, "ovs", a.Type())
}

func TestOvsAdapter_DeleteBridge(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.DeleteBridge(context.Background(), "br0")
	require.NoError(t, err)
}

func TestOvsAdapter_DeleteFlowRule(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fx := &fakeExecer{cmd: &fakeCmd{}}
		a := &OvsAdapter{ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

		err := a.DeleteFlowRule(context.Background(), "br0", "match")
		require.NoError(t, err)
	})

	t.Run("invalid bridge", func(t *testing.T) {
		fx := &fakeExecer{cmd: &fakeCmd{}}
		a := &OvsAdapter{ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

		err := a.DeleteFlowRule(context.Background(), "bad bridge", "match")
		require.Error(t, err)
		require.True(t, apperrors.Is(err, apperrors.InvalidInput))
	})
}

func TestOvsAdapter_CreateVXLANTunnel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fx := &fakeExecer{cmd: &fakeCmd{}}
		a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

		err := a.CreateVXLANTunnel(context.Background(), "br0", 100, "192.168.1.1")
		require.NoError(t, err)
	})

	t.Run("invalid bridge", func(t *testing.T) {
		fx := &fakeExecer{cmd: &fakeCmd{}}
		a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

		err := a.CreateVXLANTunnel(context.Background(), "bad bridge", 100, "192.168.1.1")
		require.Error(t, err)
		require.True(t, apperrors.Is(err, apperrors.InvalidInput))
	})
}

func TestOvsAdapter_DeleteVXLANTunnel(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.DeleteVXLANTunnel(context.Background(), "br0", "192.168.1.1")
	require.NoError(t, err)
}

func TestOvsAdapter_CreateVethPair(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{logger: slog.Default(), exec: fx}

	err := a.CreateVethPair(context.Background(), "veth0", "veth1")
	require.NoError(t, err)
}

func TestOvsAdapter_AttachVethToBridge(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{ovsPath: "/bin/ovs-vsctl", logger: slog.Default(), exec: fx}

	err := a.AttachVethToBridge(context.Background(), "br0", "veth0")
	require.NoError(t, err)
}

func TestOvsAdapter_DeleteVethPair(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{logger: slog.Default(), exec: fx}

	err := a.DeleteVethPair(context.Background(), "veth0")
	require.NoError(t, err)
}

func TestOvsAdapter_SetVethIP(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{}}
	a := &OvsAdapter{logger: slog.Default(), exec: fx}

	err := a.SetVethIP(context.Background(), "veth0", "10.0.0.1", "24")
	require.NoError(t, err)
}

func TestOvsAdapter_ListFlowRules(t *testing.T) {
	fx := &fakeExecer{cmd: &fakeCmd{out: []byte("cookie=0x0, duration=1.0s, table=0, n_packets=0, n_bytes=0, priority=100,ip actions=NORMAL\n")}}
	a := &OvsAdapter{ofctlPath: "/bin/ovs-ofctl", logger: slog.Default(), exec: fx}

	rules, err := a.ListFlowRules(context.Background(), "br0")
	require.NoError(t, err)
	require.NotNil(t, rules)
}
