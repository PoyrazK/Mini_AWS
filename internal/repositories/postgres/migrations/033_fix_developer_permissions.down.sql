-- +goose Down

-- Remove the developer permissions added in 033_fix_developer_permissions.up.sql
DELETE FROM role_permissions 
WHERE role_id = '00000000-0000-0000-0000-000000000002'
AND permission NOT IN ('instance:launch', 'instance:terminate', 'instance:read', 'vpc:read', 'volume:read', 'snapshot:read');
