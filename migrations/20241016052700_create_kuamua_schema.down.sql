
BEGIN;

-- Drop the trigger first
DROP TRIGGER IF EXISTS update_kuamua_patterns_cache_trigger ON patterns;

-- Drop the function
DROP FUNCTION IF EXISTS notify_update_kuamua_patterns_cache();

-- Drop the constraints
ALTER TABLE kuamua_patterns DROP CONSTRAINT IF EXISTS km_owneridlen;
ALTER TABLE kuamua_patterns DROP CONSTRAINT IF EXISTS km_namelen;
ALTER TABLE kuamua_patterns DROP CONSTRAINT IF EXISTS km_patternlen;
ALTER TABLE kuamua_patterns DROP CONSTRAINT IF EXISTS km_grouplen;

-- Drop the indexes
DROP INDEX IF EXISTS km_unique_name_idx;
DROP INDEX IF EXISTS km_owner_idx;
DROP INDEX IF EXISTS km_group_idx;
DROP INDEX IF EXISTS km_subgroup_idx;

-- Finally drop the table
DROP TABLE IF EXISTS kuamua_patterns;

COMMIT;
