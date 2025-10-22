-- +goose Up
-- +goose StatementBegin
INSERT INTO "public"."testing_transaction" (
    "order_no", "customer_name", "customer_email", "amount", "status"
) VALUES
    ('TEST-001', 'John Doe', 'john@example.com', 100000.00, 'pending'),
    ('TEST-002', 'Jane Smith', 'jane@example.com', 250000.00, 'paid');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."testing_transaction" WHERE "order_no" IN ('TEST-001', 'TEST-002');
-- +goose StatementEnd
