-- add token to users
ALTER TABLE "users" ADD COLUMN "token" varchar NOT NULL DEFAULT '';