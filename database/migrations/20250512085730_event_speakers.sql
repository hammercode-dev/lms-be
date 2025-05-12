-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS event_speakers_id_seq;

CREATE TABLE IF NOT EXISTS "public"."event_speakers" (
    "id" int4 NOT NULL DEFAULT nextval('event_speakers_id_seq'::regclass),
    "event_id" int4,
    "name" varchar(255),
    CONSTRAINT "event_speakers_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "public"."events"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."event_speakers";
DROP SEQUENCE IF EXISTS event_speakers_id_seq;