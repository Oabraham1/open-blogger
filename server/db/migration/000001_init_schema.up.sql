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
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "posts" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "title" varchar NOT NULL,
  "body" text NOT NULL,
  "user_id" uuid NOT NULL,
  "username" varchar NOT NULL,
  "status" status NOT NULL,
  "category" varchar NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "published_at" timestamptz NOT NULL,
  "last_modified" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "comments" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" uuid NOT NULL,
  "username" varchar NOT NULL,
  "post_id" uuid NOT NULL,
  "body" text NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

COMMENT ON COLUMN "posts"."body" IS 'Content of the blog post';

COMMENT ON COLUMN "comments"."body" IS 'Content of the comment';

ALTER TABLE "posts" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "comments" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "posts" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("username") REFERENCES "users" ("username");

ALTER TABLE "comments" ADD FOREIGN KEY ("post_id") REFERENCES "posts" ("id");
