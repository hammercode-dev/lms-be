-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS users_id_seq;

CREATE TABLE IF NOT EXISTS "public"."users" (
    "id" int8 NOT NULL DEFAULT nextval('users_id_seq'::regclass),
    "username" varchar(255) NOT NULL,
    "email" varchar(255) NOT NULL,
    "password" varchar(255) NOT NULL,
    "role" text,
    "fullname" text,
    "date_of_birth" timestamptz,
    "gender" text,
    "phone_number" text,
    "address" text,
    "github" text,
    "linkedin" text,
    "personal_web" text,
    "created_at" timestamptz,
    "updated_at" timestamptz,
    PRIMARY KEY ("id")
);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'uni_users_email'
    ) THEN
        CREATE UNIQUE INDEX uni_users_email ON public.users USING btree (email);
    END IF;
END$$;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_indexes WHERE indexname = 'uni_users_username'
    ) THEN
        CREATE UNIQUE INDEX uni_users_username ON public.users USING btree (username);
    END IF;
END$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."users";
DROP SEQUENCE IF EXISTS users_id_seq;
DROP INDEX IF EXISTS uni_users_email;
DROP INDEX IF EXISTS uni_users_username;
-- +goose StatementEnd
