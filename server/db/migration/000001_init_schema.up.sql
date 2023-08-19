CREATE TYPE "status" AS ENUM (
  'draft',
  'published'
);

CREATE TABLE "users" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "username" varchar UNIQUE NOT NULL,
  "password" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar NOT NULL,
  "interests" varchar[] NOT NULL DEFAULT '{}'::varchar[],
  "created_at" text NOT NULL DEFAULT TO_CHAR(NOW() AT TIME ZONE 'UTC', 'YYYY/MM/DD HH12:MI:SS')
);

CREATE TABLE "posts" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "title" varchar NOT NULL,
  "body" text NOT NULL,
  "username" varchar NOT NULL,
  "status" status NOT NULL,
  "category" varchar NOT NULL,
  "created_at" text NOT NULL DEFAULT TO_CHAR(NOW() AT TIME ZONE 'UTC', 'YYYY/MM/DD HH12:MI:SS'),
  "published_at" text NOT NULL,
  "last_modified" text NOT NULL DEFAULT TO_CHAR(NOW() AT TIME ZONE 'UTC', 'YYYY/MM/DD HH12:MI:SS')
);

CREATE TABLE "comments" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "username" varchar NOT NULL,
  "post_id" uuid NOT NULL,
  "body" text NOT NULL,
  "created_at" text NOT NULL DEFAULT TO_CHAR(NOW() AT TIME ZONE 'UTC', 'YYYY/MM/DD HH12:MI:SS')
);

CREATE TABLE "sessions" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "username" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "posts"."body" IS 'Content of the blog post';

COMMENT ON COLUMN "comments"."body" IS 'Content of the comment';

ALTER TABLE "posts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "sessions" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
