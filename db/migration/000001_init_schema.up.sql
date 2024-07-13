CREATE TYPE "config_types" AS ENUM (
  'auto_shoutout_activation',
  'auto_shoutout_delay',
  'blacklist',
  'timer_card_so_activation',
  'timer_card_duration'
);

CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "user_id" varchar NOT NULL,
  "user_login" varchar NOT NULL,
  "user_name" varchar NOT NULL,
  "profile_image_url" varchar,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE TABLE "streams" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "title" varchar NOT NULL,
  "game_name" varchar NOT NULL,
  "started_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "created_by" bigint NOT NULL
);

CREATE TABLE "attendances" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "stream_id" bigint UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "attendance_members" (
  "id" bigserial PRIMARY KEY,
  "attendance_id" bigint NOT NULL,
  "username" varchar UNIQUE NOT NULL,
  "is_shouted" bool NOT NULL DEFAULT false,
  "present_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "user_configs" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "config_type" config_types NOT NULL,
  "value" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "updated_at" timestamptz
);

CREATE INDEX ON "users" ("user_login");

CREATE INDEX ON "streams" ("title");

CREATE INDEX ON "streams" ("game_name");

COMMENT ON COLUMN "users"."user_id" IS 'user_id punya twitch';

COMMENT ON COLUMN "streams"."user_id" IS 'id punya user';

COMMENT ON COLUMN "attendances"."user_id" IS 'id punya user';

COMMENT ON COLUMN "user_configs"."user_id" IS 'id punya user';

ALTER TABLE "streams" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "streams" ADD FOREIGN KEY ("created_by") REFERENCES "users" ("id");

ALTER TABLE "attendances" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "attendances" ADD FOREIGN KEY ("stream_id") REFERENCES "streams" ("id");

ALTER TABLE "attendance_members" ADD FOREIGN KEY ("attendance_id") REFERENCES "attendances" ("id");

ALTER TABLE "user_configs" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
