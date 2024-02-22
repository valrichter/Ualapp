CREATE TABLE "accounts" (
    "id" serial PRIMARY KEY, "user_id" integer NOT NULL, "balance" money NOT NULL DEFAULT 0, "currency" varchar(3) NOT NULL, "created_at" timestamptz NOT NULL DEFAULT(now())
);

CREATE TABLE "entries" (
    "id" serial PRIMARY KEY, "account_id" integer NOT NULL, "amount" money NOT NULL, "created_at" timestamptz NOT NULL DEFAULT(now())
);

CREATE TABLE "transfers" (
    "id" serial PRIMARY KEY, "from_account_id" integer NOT NULL, "to_account_id" integer NOT NULL, "amount" money NOT NULL, "created_at" timestamptz NOT NULL DEFAULT(now())
);

CREATE INDEX ON "accounts" ("user_id");

-- Es mejor ALTER TABLE "accounts" ADD CONSTRAINT "unique_user_currency" UNIQUE (user_id, currency);
-- Se dejo esta para fines de documentacion
CREATE UNIQUE INDEX ON "accounts" ("user_id", "currency");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ( "from_account_id", "to_account_id" );

COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';

ALTER TABLE "accounts"
ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");

ALTER TABLE "entries"
ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");

ALTER TABLE "transfers"
ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");