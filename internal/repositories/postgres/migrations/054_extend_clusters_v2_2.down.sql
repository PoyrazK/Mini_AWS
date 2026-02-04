ALTER TABLE clusters 
DROP COLUMN IF EXISTS ha_enabled,
DROP COLUMN IF EXISTS api_server_lb_address,
DROP COLUMN IF EXISTS job_id;

ALTER TABLE load_balancers
DROP COLUMN IF EXISTS ip;
