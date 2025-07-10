-- +goose Up
-- +goose StatementBegin
TRUNCATE TABLE "public"."event_pays" RESTART IDENTITY CASCADE;

INSERT INTO "public"."event_pays" (
    "order_no", "status", "registration_event_id", "event_id", 
    "image_proof_payment", "net_amount"
) VALUES 
    ('PAY-001', 'paid', 1, 1, 'payment_proof1.jpg', 150.00),
    ('PAY-002', 'paid', 2, 2, 'payment_proof2.jpg', 299.99),
    ('PAY-003', 'pending', 3, 3, NULL, 50.00);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."event_pays" WHERE "order_no" IN ('PAY-001', 'PAY-002', 'PAY-003');
-- +goose StatementEnd
