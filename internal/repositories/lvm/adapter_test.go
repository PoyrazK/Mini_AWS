package lvm

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// fakeExecer for testing
type fakeExecer struct {
	// Map of command name to its behavior
	commands map[string]func(args ...string) ([]byte, error)
}

func newFakeExecer() *fakeExecer {
	return &fakeExecer{
		commands: make(map[string]func(args ...string) ([]byte, error)),
	}
}

func (f *fakeExecer) Run(name string, args ...string) ([]byte, error) {
	if fn, ok := f.commands[name]; ok {
		return fn(args...)
	}
	return nil, fmt.Errorf("command not found: %s", name)
}

func (f *fakeExecer) addCommand(name string, fn func(args ...string) ([]byte, error)) {
	f.commands[name] = fn
}

func TestLvmAdapter_Type(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	if adapter.Type() != "lvm" {
		t.Errorf("expected type 'lvm', got %s", adapter.Type())
	}
}

func TestLvmAdapter_Ping(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	err := adapter.Ping(context.Background())
	if err == nil {
		// If no error, perhaps lvm is available, but unlikely in test
		t.Log("lvm ping succeeded")
	} else {
		// Expected if lvm not available
		t.Logf("lvm ping failed as expected: %v", err)
	}
}

func TestLvmAdapter_CreateVolume_InvalidName(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	_, err := adapter.CreateVolume(context.Background(), "", 10)
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLvmAdapter_DeleteVolume_InvalidName(t *testing.T) {
	adapter := NewLvmAdapter("testvg")
	err := adapter.DeleteVolume(context.Background(), "")
	if err == nil {
		t.Error("expected error for empty name")
	}
}

func TestLvmAdapter_CreateVolume_Success(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvcreate", func(args ...string) ([]byte, error) {
		assert.Equal(t, []string{"-L", "10G", "-n", "test-vol", "vg0"}, args)
		return []byte("Logical volume \"test-vol\" created"), nil
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	path, err := adapter.CreateVolume(context.Background(), "test-vol", 10)

	assert.NoError(t, err)
	assert.Equal(t, "/dev/vg0/test-vol", path)
}

func TestLvmAdapter_CreateVolume_Failure(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvcreate", func(args ...string) ([]byte, error) {
		return []byte("Volume group \"vg0\" not found"), errors.New("exit status 5")
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	_, err := adapter.CreateVolume(context.Background(), "bad-vol", 5)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create logical volume")
}

func TestLvmAdapter_DeleteVolume_Success(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvremove", func(args ...string) ([]byte, error) {
		assert.Equal(t, []string{"-f", "vg0/test-vol"}, args)
		return nil, nil
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	err := adapter.DeleteVolume(context.Background(), "test-vol")

	assert.NoError(t, err)
}

func TestLvmAdapter_CreateSnapshot_Success(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvcreate", func(args ...string) ([]byte, error) {
		expected := []string{"-s", "-n", "snap1", "-L", "1G", "/dev/vg0/data-vol"}
		assert.Equal(t, expected, args)
		return nil, nil
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	err := adapter.CreateSnapshot(context.Background(), "data-vol", "snap1")

	assert.NoError(t, err)
}

func TestLvmAdapter_RestoreSnapshot_Success(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvconvert", func(args ...string) ([]byte, error) {
		expected := []string{"--merge", "vg0/snap1"}
		assert.Equal(t, expected, args)
		return nil, nil
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	err := adapter.RestoreSnapshot(context.Background(), "data-vol", "snap1")

	assert.NoError(t, err)
}

func TestLvmAdapter_DeleteSnapshot_Success(t *testing.T) {
	fake := newFakeExecer()
	fake.addCommand("lvremove", func(args ...string) ([]byte, error) {
		expected := []string{"-f", "vg0/snap1"}
		assert.Equal(t, expected, args)
		return nil, nil
	})

	adapter := &LvmAdapter{vgName: "vg0", execer: fake}
	err := adapter.DeleteSnapshot(context.Background(), "snap1")

	assert.NoError(t, err)
}

func TestLvmAdapter_AttachVolume(t *testing.T) {
	adapter := NewLvmAdapter("vg0")
	err := adapter.AttachVolume(context.Background(), "vol1", "inst1")
	assert.NoError(t, err)
}

func TestLvmAdapter_DetachVolume(t *testing.T) {
	adapter := NewLvmAdapter("vg0")
	err := adapter.DetachVolume(context.Background(), "vol1", "inst1")
	assert.NoError(t, err)
}
