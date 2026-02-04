-- Add methods column to gateway_routes
ALTER TABLE gateway_routes ADD COLUMN IF NOT EXISTS methods TEXT[] DEFAULT '{}';

-- Remove the unique constraint on path_prefix to allow multiple methods on the same path
ALTER TABLE gateway_routes DROP CONSTRAINT IF EXISTS gateway_routes_path_prefix_key;

-- Create a composite unique constraint (path_pattern, methods)
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.table_constraints WHERE constraint_name = 'gateway_routes_pattern_methods_key') THEN
        ALTER TABLE gateway_routes ADD CONSTRAINT gateway_routes_pattern_methods_key UNIQUE (path_pattern, methods);
    END IF;
END $$;

-- Create an index to speed up route lookups by method if needed
CREATE INDEX IF NOT EXISTS idx_gateway_routes_methods ON gateway_routes USING GIN (methods);
