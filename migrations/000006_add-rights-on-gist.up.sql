ALTER TABLE gists ADD column visibility VARCHAR(20) NOT NULL DEFAULT 'public';

CREATE TABLE user_gists (
  user_id UUID ,
  gist_id UUID,
  rights VARCHAR(10) NOT NULL,
  PRIMARY KEY (user_id, gist_id)
);

ALTER TABLE user_gists ADD CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE;

ALTER TABLE user_gists ADD CONSTRAINT fk_gist_id FOREIGN KEY (gist_id) REFERENCES gists(gist_id) ON DELETE CASCADE;
