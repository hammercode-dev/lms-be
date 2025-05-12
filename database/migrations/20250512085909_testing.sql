-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS testing_id_seq;

CREATE TABLE IF NOT EXISTS "public"."testing" (
    "id" int4 NOT NULL DEFAULT nextval('testing_id_seq'::regclass)
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."testing";
DROP SEQUENCE IF EXISTS testing_id_seq;