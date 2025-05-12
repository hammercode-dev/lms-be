-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS event_tags_id_seq;

CREATE TABLE IF NOT EXISTS "public"."event_tags" (
    "id" int4 NOT NULL DEFAULT nextval('event_tags_id_seq'::regclass),
    "event_id" int4,
    "tag" varchar(255),
    CONSTRAINT "event_tags_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "public"."events"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."event_tags";
DROP SEQUENCE IF EXISTS event_tags_id_seq;