CREATE TABLE "users" (
    "id" uuid PRIMARY KEY, "email" varchar(256) UNIQUE NOT NULL, "hashed_password" varchar(256) NOT NULL, "created_at" timestamptz NOT NULL DEFAULT(now()), "updated_at" timestamptz NOT NULL DEFAULT(now())
);