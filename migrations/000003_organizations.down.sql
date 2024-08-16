-- Drop the function if it exists
DROP FUNCTION IF EXISTS assert_owner_is_member();

-- Remove the foreign key constraint from the 'gists' table
ALTER TABLE gists DROP CONSTRAINT IF EXISTS fk_org_id;

-- Drop the 'org_id' column from the 'gists' table
ALTER TABLE gists DROP COLUMN IF EXISTS org_id;

-- Drop the 'member' table if it exists
DROP TABLE IF EXISTS member;

-- Drop the 'organization' table if it exists
DROP TABLE IF EXISTS organization;
