-- +goose Up
-- +goose StatementBegin
INSERT INTO "public"."images" (
    "file_name", "file_path", "format", "content_type", "is_used", "file_size"
) VALUES 
    ('workshop_banner.jpg', '/uploads/events/workshop_banner.jpg', 'jpg', 'image/jpeg', true, 256000),
    ('conference_logo.png', '/uploads/events/conference_logo.png', 'png', 'image/png', true, 128000),
    ('hackathon_poster.jpg', '/uploads/events/hackathon_poster.jpg', 'jpg', 'image/jpeg', true, 512000),
    ('payment_proof1.jpg', '/uploads/payments/payment_proof1.jpg', 'jpg', 'image/jpeg', true, 150000),
    ('payment_proof2.jpg', '/uploads/payments/payment_proof2.jpg', 'jpg', 'image/jpeg', true, 148000);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."images" WHERE "file_name" IN ('workshop_banner.jpg', 'conference_logo.png', 'hackathon_poster.jpg', 'payment_proof1.jpg', 'payment_proof2.jpg');
-- +goose StatementEnd
