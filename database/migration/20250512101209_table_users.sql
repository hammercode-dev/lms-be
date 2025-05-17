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

CREATE UNIQUE INDEX uni_users_email ON public.users USING btree (email);
CREATE UNIQUE INDEX uni_users_username ON public.users USING btree (username);

-- +goose StatementEnd
-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."users";
DROP SEQUENCE IF EXISTS users_id_seq;
-- +goose StatementEnd
