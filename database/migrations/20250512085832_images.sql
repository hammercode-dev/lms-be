-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS images_id_seq;

CREATE TABLE IF NOT EXISTS "public"."images" (
    "id" int4 NOT NULL DEFAULT nextval('images_id_seq'::regclass),
    "file_name" varchar(255),
    "file_path" varchar(255),
    "format" varchar(50),
    "content_type" varchar(100),
    "is_used" bool,
    "file_size" int8,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp,
    PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."images";
DROP SEQUENCE IF EXISTS images_id_seq;