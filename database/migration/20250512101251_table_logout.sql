-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS logout_id_seq;
CREATE TABLE IF NOT EXISTS "public"."logout" (
    "id" int8 NOT NULL DEFAULT nextval('logout_id_seq'::regclass),
    "token" varchar(255) NOT NULL,
    "expired_at" timestamptz,
    "created_at" timestamptz,
    PRIMARY KEY ("id")
);
CREATE UNIQUE INDEX uni_logout_token ON public.logout USING btree (token);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS logout_id_seq;
DROP TABLE IF EXISTS "public"."logout";
DROP SEQUENCE IF EXISTS logout_id_seq;
-- +goose StatementEnd
