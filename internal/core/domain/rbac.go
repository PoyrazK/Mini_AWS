package domain

import (
	"github.com/google/uuid"
)

type Permission string

const (
	// Compute Permissions
	PermissionInstanceLaunch    Permission = "instance:launch"
	PermissionInstanceTerminate Permission = "instance:terminate"
	PermissionInstanceRead      Permission = "instance:read"
	PermissionInstanceUpdate    Permission = "instance:update"

	// VPC Permissions
	PermissionVpcCreate Permission = "vpc:create"
	PermissionVpcDelete Permission = "vpc:delete"
	PermissionVpcRead   Permission = "vpc:read"

	// Storage Permissions
	PermissionVolumeCreate Permission = "volume:create"
	PermissionVolumeDelete Permission = "volume:delete"
	PermissionVolumeRead   Permission = "volume:read"

	// Snapshot Permissions
	PermissionSnapshotCreate  Permission = "snapshot:create"
	PermissionSnapshotDelete  Permission = "snapshot:delete"
	PermissionSnapshotRead    Permission = "snapshot:read"
	PermissionSnapshotRestore Permission = "snapshot:restore"

	// System Permissions
	PermissionFullAccess Permission = "*"
)

type Role struct {
	ID          uuid.UUID    `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Permissions []Permission `json:"permissions"`
}

// Default Roles
const (
	RoleAdmin     = "admin"
	RoleDeveloper = "developer"
	RoleViewer    = "viewer"
)
