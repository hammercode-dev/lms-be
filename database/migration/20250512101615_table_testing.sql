-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS testing_id_seq;

CREATE TABLE IF NOT EXISTS "public"."testing" (
    "id" int4 NOT NULL DEFAULT nextval('testing_id_seq'::regclass)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."testing";
DROP SEQUENCE IF EXISTS testing_id_seq;
-- +goose StatementEnd
