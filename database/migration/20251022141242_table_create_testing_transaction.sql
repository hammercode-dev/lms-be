-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS "public"."testing_transaction" (
    "id" SERIAL PRIMARY KEY,
    "order_no" VARCHAR(255) NOT NULL UNIQUE,
    "customer_name" VARCHAR(255) NOT NULL,
    "customer_email" VARCHAR(255) NOT NULL,
    "amount" DECIMAL(15, 2) NOT NULL,
    "status" VARCHAR(50) NOT NULL,
    "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."testing_transaction";
-- +goose StatementEnd
