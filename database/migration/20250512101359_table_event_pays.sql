-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS event_pays_id_seq;

CREATE TABLE IF NOT EXISTS "public"."event_pays" (
    "id" int4 NOT NULL DEFAULT nextval('event_pays_id_seq'::regclass),
    "order_no" varchar(20),
    "status" varchar(50),
    "registration_event_id" int4,
    "event_id" int4,
    "image_proof_payment" varchar(255),
    "net_amount" numeric(10,2),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    CONSTRAINT "event_pays_event_id_fkey" FOREIGN KEY ("event_id") REFERENCES "public"."events"("id") ON DELETE CASCADE,
    CONSTRAINT "event_pays_registration_event_id_fkey" FOREIGN KEY ("registration_event_id") REFERENCES "public"."registration_events"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."event_pays";
DROP SEQUENCE IF EXISTS event_pays_id_seq;
-- +goose StatementEnd
