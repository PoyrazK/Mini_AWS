package libvirt

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"strings"
	"time"

	"github.com/digitalocean/go-libvirt"
	"github.com/google/uuid"
	"github.com/poyrazk/thecloud/internal/core/ports"
)

const defaultPoolName = "default"

type LibvirtAdapter struct {
	conn   *libvirt.Libvirt
	logger *slog.Logger
	uri    string
}

func NewLibvirtAdapter(logger *slog.Logger, uri string) (*LibvirtAdapter, error) {
	if uri == "" {
		uri = "/var/run/libvirt/libvirt-sock"
	}

	// Connect to libvirt
	// We use a dialer for the unix socket
	c, err := net.DialTimeout("unix", uri, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to dial libvirt: %w", err)
	}

	l := libvirt.New(c)
	if err := l.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to libvirt: %w", err)
	}

	return &LibvirtAdapter{
		conn:   l,
		logger: logger,
		uri:    uri,
	}, nil
}

// Ping checks if libvirt is reachable
func (a *LibvirtAdapter) Ping(ctx context.Context) error {
	// Simple check: get version
	_, err := a.conn.ConnectGetLibVersion()
	return err
}

func (a *LibvirtAdapter) CreateInstance(ctx context.Context, name, imageName string, ports []string, networkID string, volumeBinds []string, env []string, cmd []string) (string, error) {
	// 1. Prepare storage
	// We assume 'imageName' is a backing volume in the default pool.
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err != nil {
		return "", fmt.Errorf("default pool not found: %w", err)
	}

	// Create root disk for the VM
	// For simplicity, we just create an empty 10GB qcows if image not found, or clone if we knew how (omitted for brevity)
	volXML := generateVolumeXML(name+"-root", 10)

	vol, err := a.conn.StorageVolCreateXML(pool, volXML, 0)
	if err != nil {
		// Try to continue if exists? No, better fail.
		return "", fmt.Errorf("failed to create root volume: %w", err)
	}

	// Get volume path
	// We need the key or path.
	// In go-libvirt, struct is returned.
	// We can construct path: /var/lib/libvirt/images/name-root
	// Or query XML.
	// For now, assume standard path.
	diskPath := fmt.Sprintf("/var/lib/libvirt/images/%s-root", name)

	// 2. Define Domain
	// Memory: 512MB
	// CPU: 1
	if networkID == "" {
		networkID = "default"
	}

	domainXML := generateDomainXML(name, diskPath, networkID, "", 512, 1)

	dom, err := a.conn.DomainDefineXML(domainXML)
	if err != nil {
		// Clean up volume
		_ = a.conn.StorageVolDelete(vol, 0)
		return "", fmt.Errorf("failed to define domain: %w", err)
	}

	// 3. Start Domain
	if err := a.conn.DomainCreate(dom); err != nil {
		return "", fmt.Errorf("failed to start domain: %w", err)
	}

	// Determine UUID (Name is used as ID here usually, but libvirt has UUID)
	// We return Name as ID to keep standard with arguments
	return name, nil
}

func (a *LibvirtAdapter) StopInstance(ctx context.Context, id string) error {
	dom, err := a.conn.DomainLookupByName(id)
	if err != nil {
		return fmt.Errorf("domain not found: %w", err)
	}

	if err := a.conn.DomainDestroy(dom); err != nil {
		return fmt.Errorf("failed to destroy domain: %w", err)
	}
	return nil
}

func (a *LibvirtAdapter) DeleteInstance(ctx context.Context, id string) error {
	dom, err := a.conn.DomainLookupByName(id)
	if err != nil {
		return nil // Assume already gone
	}

	// Stop if running
	state, _, err := a.conn.DomainGetState(dom, 0)
	if err == nil && state == 1 { // Running
		_ = a.conn.DomainDestroy(dom)
	}

	// Undefine (remove XML)
	if err := a.conn.DomainUndefine(dom); err != nil {
		return fmt.Errorf("failed to undefine domain: %w", err)
	}

	// Try to delete root volume?
	// Name-root
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err == nil {
		volName := id + "-root"
		vol, err := a.conn.StorageVolLookupByName(pool, volName)
		if err == nil {
			_ = a.conn.StorageVolDelete(vol, 0)
		}
	}

	return nil
}

func (a *LibvirtAdapter) GetInstanceLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	// Read from standard qemu log location
	// Note: This contains QEMU output, not necessarily guest console output unless serial is redirected there.
	// To get guest console, we'd need to attach to console or read a file if defined in XML.
	// Our XML defined <console type='pty'> which goes to a PTY. Reading PTY from outside is complex.
	// We'll fall back to QEMU log for debug info.
	logPath := fmt.Sprintf("/var/log/libvirt/qemu/%s.log", id)
	f, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %w", err)
	}
	return f, nil
}

func (a *LibvirtAdapter) GetInstanceStats(ctx context.Context, id string) (io.ReadCloser, error) {
	dom, err := a.conn.DomainLookupByName(id)
	if err != nil {
		return nil, fmt.Errorf("domain not found: %w", err)
	}

	// Memory stats
	// We use standard libvirt stats. The format for ComputeBackend expects JSON.
	// We'll construct a simple JSON.
	memStats, err := a.conn.DomainMemoryStats(dom, 10, 0)
	if err != nil {
		return nil, err
	}

	var usage, limit uint64
	for _, stat := range memStats {
		if stat.Tag == 6 { // rss
			usage = stat.Val * 1024 // KB to Bytes
		}
		if stat.Tag == 5 { // actual
			limit = stat.Val * 1024
		}
	}

	json := fmt.Sprintf(`{"memory_stats":{"usage":%d,"limit":%d}}`, usage, limit)
	return io.NopCloser(strings.NewReader(json)), nil
}

func (a *LibvirtAdapter) GetInstancePort(ctx context.Context, id string, internalPort string) (int, error) {
	return 0, fmt.Errorf("port forwarding not supported in libvirt adapter")
}

func (a *LibvirtAdapter) GetInstanceIP(ctx context.Context, id string) (string, error) {
	// 1. Get Domain
	dom, err := a.conn.DomainLookupByName(id)
	if err != nil {
		return "", fmt.Errorf("domain not found: %w", err)
	}

	// 2. We need the MAC address to look up DHCP leases.
	// We can get XML desc and parse it.
	xmlDesc, err := a.conn.DomainGetXMLDesc(dom, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get domain xml: %w", err)
	}

	// Extract MAC using string parsing (dirty but lighter than XML decoder for one field)
	// Look for <mac address='XX:XX:XX:XX:XX:XX'/>
	// This assumes one interface.
	start := strings.Index(xmlDesc, "<mac address='")
	if start == -1 {
		return "", fmt.Errorf("no mac address found in xml")
	}
	start += len("<mac address='")
	end := strings.Index(xmlDesc[start:], "'/>")
	if end == -1 {
		return "", fmt.Errorf("malformed xml mac address")
	}
	mac := xmlDesc[start : start+end]

	// 3. Lookup leases in default network
	// We assume "default" network
	net, err := a.conn.NetworkLookupByName("default")
	if err != nil {
		return "", fmt.Errorf("default network not found: %w", err)
	}

	// Pass nil for mac to get all leases (simplifies type handling of OptString)
	leases, _, err := a.conn.NetworkGetDhcpLeases(net, nil, 0, 0)
	if err != nil {
		return "", fmt.Errorf("failed to get leases: %w", err)
	}

	for _, lease := range leases {
		if len(lease.Mac) > 0 && lease.Mac[0] == mac {
			return lease.Ipaddr, nil
		}
	}

	return "", fmt.Errorf("no ip lease found for %s (%s)", id, mac)
}

func (a *LibvirtAdapter) Exec(ctx context.Context, id string, cmd []string) (string, error) {
	return "", fmt.Errorf("exec not supported in libvirt adapter")
}

func (a *LibvirtAdapter) RunTask(ctx context.Context, opts ports.RunTaskOptions) (string, error) {
	// For now, we assume a base image "alpine" exists in the default pool or we fail.
	// We create a new instance with a randomized name.
	name := "task-" + uuid.New().String()[:8]

	// Create volumes for binders?
	// The ports.RunTaskOptions has Binds []string.
	// We currently only support simple host binds or ignore them for the POC.
	// To perform the snapshot task, we actually rely on bind mounts.
	// As noted before, we replaced SnapshotService logic to use CreateVolumeSnapshot,
	// so RunTask is less critical for Snapshots now, but still useful for other things.

	// Since we don't have a dynamic ISO generator linked yet, we just start the VM.
	// If the user wants to run a command, we'd need Cloud-Init.

	// 1. Create a disk for the task VM (clone alpine-base if we could, or just new one)
	// We use the same create logic as CreateInstance but force a small size
	// We assume "alpine" is the image source.
	// In CreateInstance we assume image is passed as name.
	// Here opts.Image is "alpine".

	// Check if base image exists
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err != nil {
		return "", fmt.Errorf("failed to find default pool: %w", err)
	}

	// For brevity, we just create a new empty vol.
	// Real implementation should clone base image.
	volXML := generateVolumeXML(name+"-root", 1)
	vol, err := a.conn.StorageVolCreateXML(pool, volXML, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create root volume: %w", err)
	}

	// Get path
	diskPath, err := a.conn.StorageVolGetPath(vol)
	if err != nil {
		_ = a.conn.StorageVolDelete(vol, 0)
		return "", fmt.Errorf("failed to get volume path: %w", err)
	}

	// 2. Define Domain
	// We pass empty ISO for now as we don't generate one.
	domainXML := generateDomainXML(name, diskPath, "default", "", int(opts.MemoryMB), 1)

	dom, err := a.conn.DomainDefineXML(domainXML)
	if err != nil {
		_ = a.conn.StorageVolDelete(vol, 0)
		return "", fmt.Errorf("failed to define domain: %w", err)
	}

	// 3. Start Domain
	if err := a.conn.DomainCreate(dom); err != nil {
		_ = a.conn.DomainUndefine(dom)
		_ = a.conn.StorageVolDelete(vol, 0)
		return "", fmt.Errorf("failed to start domain: %w", err)
	}

	return name, nil
}

func (a *LibvirtAdapter) WaitTask(ctx context.Context, id string) (int64, error) {
	// Poll for domain state to be Shutoff
	// Since we can't easily get the exit code from inside the VM without qemu-agent,
	// we assume 0 if it shuts down gracefully (state Shutoff).

	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return -1, ctx.Err()
		case <-ticker.C:
			dom, err := a.conn.DomainLookupByName(id)
			if err != nil {
				// If domain is gone, maybe it was deleted?
				return -1, fmt.Errorf("domain not found: %w", err)
			}

			state, _, err := a.conn.DomainGetState(dom, 0)
			if err != nil {
				continue
			}

			// libvirt.DomainShutoff = 5
			if state == 5 {
				return 0, nil
			}
		}
	}
}

func (a *LibvirtAdapter) CreateNetwork(ctx context.Context, name string) (string, error) {
	// Simple NAT network
	xml := generateNetworkXML(name, "virbr-"+name, "192.168.123.1", "192.168.123.2", "192.168.123.254")

	net, err := a.conn.NetworkDefineXML(xml)
	if err != nil {
		return "", fmt.Errorf("failed to define network: %w", err)
	}

	if err := a.conn.NetworkCreate(net); err != nil {
		return "", fmt.Errorf("failed to start network: %w", err)
	}

	return net.Name, nil
}

func (a *LibvirtAdapter) DeleteNetwork(ctx context.Context, id string) error {
	net, err := a.conn.NetworkLookupByName(id)
	if err != nil {
		return nil // assume deleted
	}

	if err := a.conn.NetworkDestroy(net); err != nil {
		a.logger.Warn("failed to destroy network", "error", err)
	}
	if err := a.conn.NetworkUndefine(net); err != nil {
		return fmt.Errorf("failed to undefine network: %w", err)
	}
	return nil
}

func (a *LibvirtAdapter) CreateVolume(ctx context.Context, name string) error {
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err != nil {
		return fmt.Errorf("failed to find default storage pool: %w", err)
	}

	// 10GB default
	xml := generateVolumeXML(name, 10)

	// Refresh pool first
	if err := a.conn.StoragePoolRefresh(pool, 0); err != nil {
		// Log but continue
		a.logger.Warn("failed to refresh pool", "error", err)
	}

	_, err = a.conn.StorageVolCreateXML(pool, xml, 0)
	if err != nil {
		return fmt.Errorf("failed to create volume xml: %w", err)
	}
	return nil
}

func (a *LibvirtAdapter) DeleteVolume(ctx context.Context, name string) error {
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err != nil {
		return fmt.Errorf("failed to find default storage pool: %w", err)
	}

	vol, err := a.conn.StorageVolLookupByName(pool, name)
	if err != nil {
		// Check if not found
		return nil
	}

	if err := a.conn.StorageVolDelete(vol, 0); err != nil {
		return fmt.Errorf("failed to delete volume: %w", err)
	}
	return nil
}

func (a *LibvirtAdapter) CreateVolumeSnapshot(ctx context.Context, volumeID string, destinationPath string) error {
	// volumeID is the libvirt volume name
	pool, err := a.conn.StoragePoolLookupByName(defaultPoolName)
	if err != nil {
		return fmt.Errorf("failed to find default storage pool: %w", err)
	}

	vol, err := a.conn.StorageVolLookupByName(pool, volumeID)
	if err != nil {
		return fmt.Errorf("failed to find volume: %w", err)
	}
	_ = vol // Prevent unused variable error

	// We can read volume content via stream or direct file access if local.
	// For simulation on local fs, we can just copy the file?
	// But qcow2 is a format. We want a tarball of the FILESYSTEM inside the qcow2?
	// SnapshotService expects a tarball of the content.
	// Opening a qcow2 and mounting it requires nbd or guestmount.
	// This is getting complex for a "Mini AWS" without root.

	// If we assume the "volume" is just a raw file we can tar it.
	// But we initialized it as qcow2.

	// For now, return not implemented to allow compilation but indicate gap.
	return fmt.Errorf("not implemented: libvirt volume snapshot requires running agent or qemu-img convert")
}

func (a *LibvirtAdapter) RestoreVolumeSnapshot(ctx context.Context, volumeID string, sourcePath string) error {
	return fmt.Errorf("not implemented: libvirt volume restore")
}
