-- +migrate Up
CREATE SEQUENCE IF NOT EXISTS registration_events_id_seq;

CREATE TABLE IF NOT EXISTS "public"."registration_events" (
    "id" int4 NOT NULL DEFAULT nextval('registration_events_id_seq'::regclass),
    "order_no" varchar(255),
    "event_id" int4,
    "user_id" varchar(255),
    "name" varchar(255),
    "email" varchar(255),
    "phone_number" varchar(255),
    "image_proof_payment" varchar(255),
    "payment_date" timestamp,
    "status" varchar(50),
    "up_to_you" varchar(255),
    "created_by_user_id" int8,
    "updated_by_user_id" int8,
    "deleted_by_user_id" int8,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    "deleted_at" timestamp,
    CONSTRAINT "registration_events_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "public"."events"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);

-- +migrate Down
DROP TABLE IF EXISTS "public"."registration_events";
DROP SEQUENCE IF EXISTS registration_events_id_seq;