-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS transaction_events_id_seq;

CREATE TABLE IF NOT EXISTS "public"."transaction_events" (
    "id" int4 NOT NULL DEFAULT nextval('transaction_events_id_seq'::regclass),
    "transaction_no" varchar(100) UNIQUE NOT NULL,
    "registration_id" int4 NOT NULL,
    "amount" numeric(15,2) NOT NULL,
    "status" varchar(50) NOT NULL DEFAULT 'pending',
    "invoice_id" varchar(255),
    "invoice_url" text,
    "external_id" varchar(255) UNIQUE,
    "payment_method" varchar(50),
    "paid_at" timestamp,
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp,
    PRIMARY KEY ("id"),
    CONSTRAINT "fk_transaction_registration" 
        FOREIGN KEY ("registration_id") 
        REFERENCES "public"."registration_events"("id") 
        ON DELETE CASCADE
);

-- Create indexes for performance
CREATE INDEX idx_trx_no ON transaction_events(transaction_no);
CREATE INDEX idx_trx_registration ON transaction_events(registration_id);
CREATE INDEX idx_trx_status ON transaction_events(status);
CREATE INDEX idx_trx_external_id ON transaction_events(external_id);

COMMENT ON TABLE "public"."transaction_events" IS 'Payment transactions for event registrations';
COMMENT ON COLUMN "public"."transaction_events"."status" IS 'pending, paid, expired';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."transaction_events";
DROP SEQUENCE IF EXISTS transaction_events_id_seq;
-- +goose StatementEnd
