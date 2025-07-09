-- +goose Up
-- +goose StatementBegin
TRUNCATE TABLE "public"."events" RESTART IDENTITY CASCADE;

INSERT INTO "public"."events" (
    "id", "title", "description", "author", "image", "date",
    "reservation_start_date", "reservation_end_date", "type", 
    "location", "duration", "status", "price", "capacity",
    "registration_link", "created_at", "updated_at"
) VALUES 
    (1, 'Web Development Workshop', 'Learn modern web development techniques', 
    'Tech Team', 'workshop_banner.jpg', '2025-06-15 09:00:00',
    '2025-05-15 00:00:00', '2025-06-14 23:59:59', 'workshop',
    'Tech Hub, Floor 3', '8 hours', 'upcoming', 150.00, 30,
    'https://example.com/register/web-dev', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    (2, 'Data Science Conference', 'Annual conference for data scientists', 
    'Data Science Society', 'conference_logo.png', '2025-07-20 08:30:00',
    '2025-06-01 00:00:00', '2025-07-15 23:59:59', 'conference',
    'Convention Center', '3 days', 'upcoming', 299.99, 200,
    'https://example.com/register/data-conf', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP),
    
    (3, 'Mobile App Hackathon', '48-hour app building competition', 
    'Mobile Developers Group', 'hackathon_poster.jpg', '2025-08-10 10:00:00',
    '2025-07-01 00:00:00', '2025-08-05 23:59:59', 'hackathon',
    'Innovation Labs', '48 hours', 'upcoming', 50.00, 100,
    'https://example.com/register/app-hackathon', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."events" WHERE "id" IN (1, 2, 3);
-- +goose StatementEnd
