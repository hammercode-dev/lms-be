-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='blog_posts' AND column_name='created_at') THEN
        ALTER TABLE blog_posts ADD COLUMN created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='blog_posts' AND column_name='is_deleted') THEN
        ALTER TABLE blog_posts ADD COLUMN is_deleted BOOLEAN DEFAULT FALSE;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
ALTER TABLE blog_posts DROP COLUMN IF EXISTS created_at;
ALTER TABLE blog_posts DROP COLUMN IF EXISTS is_deleted;
