-- +goose Up

-- Add missing cluster permissions for developer role
INSERT INTO role_permissions (role_id, permission) VALUES 
    ('00000000-0000-0000-0000-000000000002', 'cluster:create'),
    ('00000000-0000-0000-0000-000000000002', 'cluster:read'),
    ('00000000-0000-0000-0000-000000000002', 'cluster:update'),
    ('00000000-0000-0000-0000-000000000002', 'cluster:delete'),
    ('00000000-0000-0000-0000-000000000002', 'cluster:list'),
    ('00000000-0000-0000-0000-000000000002', 'security_group:create'),
    ('00000000-0000-0000-0000-000000000002', 'security_group:read'),
    ('00000000-0000-0000-0000-000000000002', 'security_group:update'),
    ('00000000-0000-0000-0000-000000000002', 'security_group:delete')
ON CONFLICT (role_id, permission) DO NOTHING;
