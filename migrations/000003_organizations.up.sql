CREATE TABLE IF NOT EXISTS organization (
    org_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS member (
    member_id SERIAL PRIMARY KEY,
    org_id INT NOT NULL,
    user_id INT NOT NULL,
    CONSTRAINT fk_org_id FOREIGN KEY (org_id) REFERENCES organization(org_id) ON DELETE CASCADE,
    CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

ALTER TABLE gists ADD COLUMN org_id INT NULL;
ALTER TABLE gists ADD CONSTRAINT fk_org_id FOREIGN KEY (org_id) REFERENCES organization(org_id) ON DELETE CASCADE;

CREATE OR REPLACE FUNCTION assert_owner_is_member() RETURNS TRIGGER AS $$
BEGIN
    IF (SELECT COUNT(*) FROM member WHERE org_id = NEW.org_id AND user_id = NEW.owner) = 0 THEN
        RAISE EXCEPTION 'Owner is not a member of the organization';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER assert_owner_is_member
BEFORE INSERT ON gists
FOR EACH ROW
EXECUTE FUNCTION assert_owner_is_member();