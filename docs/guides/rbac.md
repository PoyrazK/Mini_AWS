# Role-Based Access Control (RBAC)

This guide explains the RBAC system: roles, permissions, APIs, and CLI commands for management.

## Overview

RBAC provides a way to assign permissions to users via roles. The system supports:

- **Roles**: Collections of permissions (e.g., `admin`, `developer`, `viewer`)
- **Permissions**: Fine-grained actions on resources (e.g., `instance:launch`, `volume:snapshot`)
- **Role bindings**: Assign a role to a user (via Email or User ID)

## CLI Management

The `cloud roles` command-set provides full management capabilities.

### Create a Role
```bash
cloud roles create developer --description "Developer access" --permissions "instance:read,instance:launch,volume:read"
```

### List Roles
```bash
cloud roles list
```

### Bind Role to User
You can bind a role using either the user's ID or their email address.
```bash
cloud roles bind user@example.com developer
```

### List Role Bindings
See which users have which roles assigned.
```bash
cloud roles list-bindings
```

### Delete a Role
```bash
cloud roles delete <role-uuid>
```

## Default Roles (Fallback)

The system includes hardcoded fallbacks for default roles if they are not yet defined in the database:

- **admin**: Has `*` (Full Access) permission.
- **developer**: Has all permissions except RBAC management.
- **viewer**: Has read-only permissions for instances, volumes, VPCs, and snapshots.

## API Reference

### Roles
- `POST /rbac/roles` — Create a new role
- `GET /rbac/roles` — List all roles
- `GET /rbac/roles/:id` — Get role details
- `PUT /rbac/roles/:id` — Update role name, description, or permissions
- `DELETE /rbac/roles/:id` — Delete a role

### Permissions
- `POST /rbac/roles/:id/permissions` — Add a single permission to a role
- `DELETE /rbac/roles/:id/permissions/:permission` — Remove a permission from a role

### Bindings
- `POST /rbac/bindings` — Create a role binding (`user_identifier` and `role_name`)
- `GET /rbac/bindings` — List all role bindings (returns user list with roles)

## Permission List

Typical permissions include:
- `*` (Full Access)
- `instance:launch`, `instance:read`, `instance:update`, `instance:terminate`
- `volume:create`, `volume:read`, `volume:delete`, `volume:snapshot`
- `vpc:create`, `vpc:read`, `vpc:delete`
- `rbac:manage`

## Troubleshooting

- **403 Forbidden**: Ensure your user has the required permission or the `*` wildcard.
- **Role Not Found**: Check if the role name matches exactly (case-sensitive).
- **Binding Failed**: Ensure the user email exists in the system.
