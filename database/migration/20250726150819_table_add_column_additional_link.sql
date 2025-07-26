-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE events
ADD COLUMN additional_link TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE events
DROP COLUMN additional_link;
-- +goose StatementEnd
