-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."testing_transaction" (
    "id" SERIAL PRIMARY KEY,
    "order_no" varchar(100) UNIQUE NOT NULL,
    "customer_name" varchar(255) NOT NULL,
    "customer_email" varchar(255) NOT NULL,
    "amount" numeric(15,2) NOT NULL,
    "status" varchar(50) NOT NULL DEFAULT 'pending',
    "invoice_url" text,
    "payment_method" varchar(50),
    "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp
);

COMMENT ON TABLE "public"."testing_transaction" IS 'Simple table for testing Xendit payment gateway';
COMMENT ON COLUMN "public"."testing_transaction"."status" IS 'pending, paid, expired';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."testing_transaction";
-- +goose StatementEnd
