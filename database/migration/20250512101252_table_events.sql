-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS events_id_seq;

CREATE TABLE IF NOT EXISTS "public"."events" (
    "id" int4 NOT NULL DEFAULT nextval('events_id_seq'::regclass),
    "title" varchar(255),
    "description" text,
    "author" varchar(255),
    "image" varchar(255),
    "date" timestamp,
    "reservation_start_date" timestamp,
    "reservation_end_date" timestamp,
    "type" varchar(50),
    "location" varchar(255),
    "duration" varchar(50),
    "status" varchar(50),
    "price" numeric(10,2),
    "capacity" int4,
    "registration_link" varchar(255),
    "created_by" int4,
    "updated_by" int4,
    "deleted_by" int4,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    "is_online" bool,
    "slug" varchar,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."events";
DROP SEQUENCE IF EXISTS events_id_seq;
-- +goose StatementEnd
