-- add column encrypted_twitch_token
ALTER TABLE sessions ADD COLUMN encrypted_twitch_token VARCHAR NOT NULL DEFAULT '';