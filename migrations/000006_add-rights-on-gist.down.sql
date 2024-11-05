ALTER TABLE user_gists DROP CONSTRAINT fk_user_id;
ALTER TABLE user_gists DROP CONSTRAINT fk_gist_id;

DROP TABLE user_gists;

ALTER TABLE gists DROP COLUMN visibility;
