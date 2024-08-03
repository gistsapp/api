CREATE TABLE IF NOT EXISTS auth_identity (
  auth_id SERIAL PRIMARY KEY,
  data JSONB NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  user_id SERIAL PRIMARY KEY,
  email VARCHAR(320) NOT NULL,
  name VARCHAR(300) NOT NULL
);

CREATE TABLE IF NOT EXISTS users_auth (
  user_id INT NOT NULL,
  auth_id INT NOT NULL,
  PRIMARY KEY (user_id, auth_id),
  CONSTRAINT fk_user_id FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
  CONSTRAINT fk_auth_id FOREIGN KEY (auth_id) REFERENCES auth_identity(auth_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS token (
  token_id SERIAL PRIMARY KEY,
  type VARCHAR(10) NOT NULL,
  value VARCHAR(200) NOT NULL
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
