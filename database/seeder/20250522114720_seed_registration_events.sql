-- +goose Up
-- +goose StatementBegin
INSERT INTO "public"."registration_events" (
    "order_no", "event_id", "user_id", "name", "email", "phone_number", 
    "image_proof_payment", "payment_date", "status", "up_to_you", "created_by_user_id"
) VALUES 
    ('ORD-001', 1, '2', 'John Doe', 'john@example.com', '987654321', 
     'payment_proof1.jpg', '2025-06-01 14:30:00', 'confirmed', 'Looking forward to it!', 2),
    
    ('ORD-002', 2, '3', 'Jane Doe', 'jane@example.com', '555123456', 
     'payment_proof2.jpg', '2025-07-05 10:15:00', 'confirmed', 'Excited to learn!', 3),
     
    ('ORD-003', 3, '2', 'John Doe', 'john@example.com', '987654321', 
     NULL, NULL, 'pending', 'Will bring my team', 2);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."registration_events" WHERE "order_no" IN ('ORD-001', 'ORD-002', 'ORD-003');
-- +goose StatementEnd
