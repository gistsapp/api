CREATE TABLE IF NOT EXISTS users (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(320) NOT NULL,
  name VARCHAR(300) NOT NULL,
  picture TEXT
);

CREATE TABLE IF NOT EXISTS token (
  token_id SERIAL PRIMARY KEY,
  type VARCHAR(10) NOT NULL,
  value VARCHAR(200) NOT NULL,
  keyword TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ DEFAULT NOW()
);


CREATE INDEX IF NOT EXISTS idx_token_value ON token(keyword);

CREATE OR REPLACE FUNCTION delete_old_kv_pairs() RETURNS TRIGGER AS $$
BEGIN
    DELETE FROM token WHERE NEW.created_at < NOW() - INTERVAL '30 minutes';
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER kv_store_cleanup
BEFORE INSERT ON token
FOR EACH ROW
EXECUTE FUNCTION delete_old_kv_pairs();


CREATE TABLE IF NOT EXISTS auth_identity (
  auth_id SERIAL PRIMARY KEY,
  provider_id TEXT NOT NULL,
  data JSONB NOT NULL,
  type VARCHAR(10) NOT NULL,
  owner_id INT NOT NULL,
  CONSTRAINT fk_owner_id FOREIGN KEY (owner_id) REFERENCES users(user_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users_token (
  user_id INT NOT NULL,
  token_id INT NOT NULL,
  PRIMARY KEY (user_id, token_id),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  CONSTRAINT fk_token_id FOREIGN KEY (token_id) REFERENCES token(token_id) ON DELETE CASCADE
);

ALTER TABLE IF EXISTS gists ADD COLUMN owner INT NOT NULL;
ALTER TABLE IF EXISTS gists ADD CONSTRAINT fk_owner FOREIGN KEY (owner) REFERENCES users(user_id) ON DELETE CASCADE;
