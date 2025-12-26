-- +goose Up
-- +goose StatementBegin
Truncate Table "public"."event_tags" Restart Identity Cascade;
INSERT INTO "public"."event_tags" (
    "event_id", "tag"
) VALUES 
    (1, 'programming'),
    (1, 'web-development'),
    (1, 'javascript'),
    (2, 'data-science'),
    (2, 'machine-learning'),
    (2, 'big-data'),
    (3, 'mobile'),
    (3, 'hackathon'),
    (3, 'app-development');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."event_tags" WHERE "event_id" IN (1, 2, 3);
-- +goose StatementEnd
