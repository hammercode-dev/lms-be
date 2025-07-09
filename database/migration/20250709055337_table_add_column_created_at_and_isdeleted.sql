-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE blog_posts
  ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;



-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blog_posts CASCADE;
-- +goose StatementEnd
