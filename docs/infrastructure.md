# Cloud Infrastructure Guide

This document covers the infrastructure and DevOps aspects of The Cloud.

## Compute Backends (`internal/repositories`)

The Cloud supports multiple compute backends, selectable via the `COMPUTE_BACKEND` environment variable.

### 1. Docker Adapter (`/docker`)
Launches **Docker Containers** that act as instances.
- **Pull**: Automatically pulls images from Docker Hub.
- **Isolation**: Shared kernel, rapid startup.
- **Limitation**: No block device attachment support (mapped to bind mounts).

### 2. Libvirt Adapter (`/libvirt`)
Launches full **KVM Virtual Machines**.
- **Template-based**: Uses XML templates to define VM hardware.
- **VNC Console**: Supports remote console access via VNC with dynamic port assignment.
- **Disk Support**: Handles both file-based QCOW2 images and raw LVM block devices.
- **Requirements**: Requires `libvirtd` running on the host and a configured storage pool.

## Storage Backends (`internal/repositories`)

The `StorageBackend` port abstracts volume management from the compute lifecycle.

### 1. LVM Adapter (`/lvm`)
Directly manages hardware block storage.
- **Volume Creation**: `lvcreate` is used to allocate exact GB segments from a Volume Group.
- **Snapshots**: Uses LVM's native COW snapshots (`lvcreate -s`).
- **Restoration**: Fast merge-back restoration using `lvconvert --merge`.
- **Requirements**: A Volume Group named `thecloud-vg` (configurable).

### 2. Noop Adapter (`/noop`)
Used for development and testing environments without physical storage hardware.

## Network Architecture (SDN)

The Cloud uses **Open vSwitch (OVS)** for Software Defined Networking.

- **VPC Isolation**: Each VPC is mapped to a logical bridge or VXLAN tunnel.
- **Flow Control**: OVS flow rules enforce isolation between different users' subnets.
- **Veth Plumbing**: Instances are connected to OVS bridges via veth pairs.

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | API Server Port | `8080` |
| `DATABASE_URL` | Postgres Connection String | `postgres://localhost...` |
| `COMPUTE_BACKEND` | `docker`, `libvirt`, or `noop` | `docker` |
| `STORAGE_BACKEND` | `lvm` or `noop` | `noop` |
| `NETWORK_BACKEND` | `ovs` or `noop` | `ovs` |

## Deployment Strategy

### Docker Compose
Used to run the control plane services (API, Postgres, Redis).
```yaml
services:
  api:
    environment:
      - COMPUTE_BACKEND=docker
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
```

### Bare Metal / KVM Host
For production-like Libvirt/LVM setups, the API should run on a host with:
1. `libvirtd` installed and running.
2. `lvm2` tools installed.
3. `openvswitch-switch` installed.
