ALTER TABLE member DROP CONSTRAINT fk_org_id; 
ALTER TABLE member DROP CONSTRAINT fk_user_id; 

ALTER TABLE auth_identity DROP CONSTRAINT fk_owner_id; 

ALTER TABLE gists DROP CONSTRAINT fk_owner; 
ALTER TABLE gists DROP CONSTRAINT fk_org_id;

ALTER TABLE auth_identity DROP CONSTRAINT auth_identity_pkey; 
ALTER TABLE gists DROP CONSTRAINT gists_pkey; 
ALTER TABLE member DROP CONSTRAINT member_pkey;
ALTER TABLE organization DROP CONSTRAINT organization_pkey;
ALTER TABLE token DROP CONSTRAINT token_pkey;
ALTER TABLE users DROP CONSTRAINT users_pkey;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

ALTER TABLE users ADD COLUMN user_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE member ADD COLUMN user_id_new UUID;
ALTER TABLE gists ADD COLUMN owner_new UUID;
ALTER TABLE auth_identity ADD COLUMN owner_id_new UUID;

CREATE OR REPLACE PROCEDURE update_users_id()
LANGUAGE plpgsql
AS $$
DECLARE
  v_user RECORD;
BEGIN
  UPDATE users SET user_id_new = uuid_generate_v4();

  FOR v_user IN (SELECT user_id, user_id_new FROM users) LOOP
    UPDATE member
    SET user_id_new = v_user.user_id_new
    WHERE user_id = v_user.user_id;
    
    UPDATE gists
    SET owner_new = v_user.user_id_new
    WHERE owner = v_user.user_id;

    UPDATE auth_identity
    SET owner_id_new = v_user.user_id_new
    WHERE owner_id = v_user.user_id;
  END LOOP;
  

END;
$$;

CALL update_users_id();


ALTER TABLE gists DROP COLUMN owner;
ALTER TABLE gists RENAME COLUMN owner_new TO owner;

ALTER TABLE member DROP COLUMN user_id;
ALTER TABLE member RENAME COLUMN user_id_new TO user_id;

ALTER TABLE auth_identity DROP COLUMN owner_id;
ALTER TABLE auth_identity RENAME COLUMN owner_id_new TO owner_id;

ALTER TABLE users DROP COLUMN user_id;
ALTER TABLE users RENAME COLUMN user_id_new TO user_id;

ALTER TABLE users ADD PRIMARY KEY (user_id);
ALTER TABLE member ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE;
ALTER TABLE auth_identity ADD CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE;
ALTER TABLE gists ADD CONSTRAINT fk_owner FOREIGN KEY (owner) REFERENCES users(user_id) ON DELETE CASCADE;


ALTER TABLE organization ADD COLUMN org_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE member ADD COLUMN org_id_new UUID;
ALTER TABLE gists ADD COLUMN org_id_new UUID;

CREATE OR REPLACE PROCEDURE update_org_id()
LANGUAGE plpgsql
AS $$
DECLARE
  v_org RECORD;
BEGIN
  UPDATE organization SET org_id_new = uuid_generate_v4();

  FOR v_org IN (SELECT org_id, org_id_new FROM organization) LOOP
    UPDATE member
    SET org_id_new = v_org.org_id_new
    WHERE org_id = v_org.org_id;
    
    UPDATE gists
    SET org_id_new = v_org.org_id_new
    WHERE org_id = v_org.org_id;

  END LOOP;
END;
$$;

CALL update_org_id();

ALTER TABLE gists DROP COLUMN org_id;
ALTER TABLE gists RENAME COLUMN org_id_new TO org_id;

ALTER TABLE member DROP COLUMN org_id;
ALTER TABLE member RENAME COLUMN org_id_new TO org_id;

ALTER TABLE organization DROP COLUMN org_id;
ALTER TABLE organization RENAME COLUMN org_id_new TO org_id;
ALTER TABLE organization ADD PRIMARY KEY (org_id);
ALTER TABLE member ADD CONSTRAINT fk_org_id FOREIGN KEY (org_id) REFERENCES organization(org_id) ON DELETE CASCADE;

ALTER TABLE gists ADD CONSTRAINT fk_org_id FOREIGN KEY (org_id) REFERENCES organization(org_id) ON DELETE CASCADE;


ALTER TABLE auth_identity ADD COLUMN auth_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE auth_identity DROP COLUMN auth_id;
ALTER TABLE auth_identity RENAME COLUMN auth_id_new TO auth_id;
ALTER TABLE auth_identity ADD PRIMARY KEY (auth_id);

ALTER TABLE gists ADD COLUMN gist_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE gists DROP COLUMN gist_id;
ALTER TABLE gists RENAME COLUMN gist_id_new TO gist_id;
ALTER TABLE gists ADD PRIMARY KEY (gist_id);

ALTER TABLE member ADD COLUMN member_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE member DROP COLUMN member_id;
ALTER TABLE member RENAME COLUMN member_id_new TO member_id;
ALTER TABLE member ADD PRIMARY KEY (member_id);

ALTER TABLE token ADD COLUMN token_id_new UUID DEFAULT uuid_generate_v4();
ALTER TABLE token DROP COLUMN token_id;
ALTER TABLE token RENAME COLUMN token_id_new TO token_id;
ALTER TABLE token ADD PRIMARY KEY (token_id);
