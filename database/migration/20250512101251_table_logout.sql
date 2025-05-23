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

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'uni_logout_token'
    ) THEN
        CREATE UNIQUE INDEX uni_logout_token ON public.logout USING btree (token);
    END IF;
END$$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS uni_logout_token;
DROP TABLE IF EXISTS "public"."logout";
DROP SEQUENCE IF EXISTS logout_id_seq;
-- +goose StatementEnd
