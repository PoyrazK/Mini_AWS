-- 1. Create personal tenants for existing users
INSERT INTO tenants (id, name, slug, owner_id, plan, status, created_at, updated_at)
SELECT 
    gen_random_uuid(), 
    name || '''s Personal Tenant', 
    'personal-' || LOWER(REPLACE(name, ' ', '-')) || '-' || SUBSTR(id::text, 1, 8),
    id, 
    'free', 
    'active', 
    NOW(), 
    NOW()
FROM users;

-- 2. Add users as members of their new tenants
INSERT INTO tenant_members (tenant_id, user_id, role, joined_at)
SELECT id, owner_id, 'owner', NOW()
FROM tenants;

-- 3. Update users table with default_tenant_id
UPDATE users u
SET default_tenant_id = t.id
FROM tenants t
WHERE u.id = t.owner_id;

-- 4. Update all resources to match their owner's personal tenant
UPDATE instances i SET tenant_id = u.default_tenant_id FROM users u WHERE i.user_id = u.id;
UPDATE vpcs v SET tenant_id = u.default_tenant_id FROM users u WHERE v.user_id = u.id;
UPDATE volumes vol SET tenant_id = u.default_tenant_id FROM users u WHERE vol.user_id = u.id;
UPDATE load_balancers lb SET tenant_id = u.default_tenant_id FROM users u WHERE lb.user_id = u.id;
UPDATE scaling_groups sg SET tenant_id = u.default_tenant_id FROM users u WHERE sg.user_id = u.id;
UPDATE buckets b SET tenant_id = u.default_tenant_id FROM users u WHERE b.user_id = u.id;
UPDATE clusters c SET tenant_id = u.default_tenant_id FROM users u WHERE c.user_id = u.id;
UPDATE floating_ips f SET tenant_id = u.default_tenant_id FROM users u WHERE f.user_id = u.id;
UPDATE security_groups s SET tenant_id = u.default_tenant_id FROM users u WHERE s.user_id = u.id;
UPDATE subnets sub SET tenant_id = u.default_tenant_id FROM users u WHERE sub.user_id = u.id;
UPDATE snapshots snap SET tenant_id = u.default_tenant_id FROM users u WHERE snap.user_id = u.id;
UPDATE secrets sec SET tenant_id = u.default_tenant_id FROM users u WHERE sec.user_id = u.id;
UPDATE functions func SET tenant_id = u.default_tenant_id FROM users u WHERE func.user_id = u.id;
UPDATE caches c SET tenant_id = u.default_tenant_id FROM users u WHERE c.user_id = u.id;
UPDATE queues q SET tenant_id = u.default_tenant_id FROM users u WHERE q.user_id = u.id;
UPDATE topics t SET tenant_id = u.default_tenant_id FROM users u WHERE t.user_id = u.id;
UPDATE subscriptions s SET tenant_id = u.default_tenant_id FROM users u WHERE s.user_id = u.id;
UPDATE cron_jobs c SET tenant_id = u.default_tenant_id FROM users u WHERE c.user_id = u.id;
UPDATE gateway_routes g SET tenant_id = u.default_tenant_id FROM users u WHERE g.user_id = u.id;
UPDATE deployments d SET tenant_id = u.default_tenant_id FROM users u WHERE d.user_id = u.id;
UPDATE stacks s SET tenant_id = u.default_tenant_id FROM users u WHERE s.user_id = u.id;
UPDATE images img SET tenant_id = u.default_tenant_id FROM users u WHERE img.user_id = u.id;

-- 5. Add default quotas for all tenants
INSERT INTO tenant_quotas (tenant_id, max_instances, max_vpcs, max_storage_gb, max_memory_gb, max_vcpus)
SELECT id, 10, 3, 100, 32, 16 FROM tenants;

-- 6. Update API keys with default_tenant_id
UPDATE api_keys ak
SET default_tenant_id = u.default_tenant_id
FROM users u
WHERE ak.user_id = u.id;
