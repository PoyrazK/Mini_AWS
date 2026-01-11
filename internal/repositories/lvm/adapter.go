package lvm

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/poyrazk/thecloud/internal/core/ports"
)

type LvmAdapter struct {
	vgName string
}

func NewLvmAdapter(vgName string) *LvmAdapter {
	return &LvmAdapter{vgName: vgName}
}

func (a *LvmAdapter) CreateVolume(ctx context.Context, name string, sizeGB int) (string, error) {
	// lvcreate -L 10G -n vol_name vg_name
	cmd := exec.CommandContext(ctx, "lvcreate", "-L", fmt.Sprintf("%dG", sizeGB), "-n", name, a.vgName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to create logical volume: %v, output: %s", err, string(out))
	}
	return fmt.Sprintf("/dev/%s/%s", a.vgName, name), nil
}

func (a *LvmAdapter) DeleteVolume(ctx context.Context, name string) error {
	// lvremove -f vg_name/vol_name
	cmd := exec.CommandContext(ctx, "lvremove", "-f", fmt.Sprintf("%s/%s", a.vgName, name))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove logical volume: %v, output: %s", err, string(out))
	}
	return nil
}

func (a *LvmAdapter) AttachVolume(ctx context.Context, volumeName, instanceID string) error {
	// Attaching in LVM context usually means making it available to the hypervisor.
	// For Libvirt, it's about adding the disk to the XML.
	// This might be better handled in the Compute Service or by a direct Libvirt call.
	// For now, we'll consider it a no-op if the volume is already in /dev/vg/vol.
	return nil
}

func (a *LvmAdapter) DetachVolume(ctx context.Context, volumeName, instanceID string) error {
	return nil
}

func (a *LvmAdapter) CreateSnapshot(ctx context.Context, volumeName, snapshotName string) error {
	// lvcreate -s -n snapshot_name -L 1G /dev/vg_name/volume_name
	// Note: LVM snapshots need space. For simplicity, we use a fixed size or same as original.
	cmd := exec.CommandContext(ctx, "lvcreate", "-s", "-n", snapshotName, "-L", "1G", fmt.Sprintf("/dev/%s/%s", a.vgName, volumeName))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to create lvm snapshot: %v, output: %s", err, string(out))
	}
	return nil
}

func (a *LvmAdapter) DeleteSnapshot(ctx context.Context, snapshotName string) error {
	cmd := exec.CommandContext(ctx, "lvremove", "-f", fmt.Sprintf("%s/%s", a.vgName, snapshotName))
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove lvm snapshot: %v, output: %s", err, string(out))
	}
	return nil
}

func (a *LvmAdapter) Ping(ctx context.Context) error {
	// Check if vgs command works and vg exists
	cmd := exec.CommandContext(ctx, "vgs", a.vgName)
	return cmd.Run()
}

func (a *LvmAdapter) Type() string {
	return "lvm"
}

// Ensure LvmAdapter implements StorageBackend
var _ ports.StorageBackend = (*LvmAdapter)(nil)
