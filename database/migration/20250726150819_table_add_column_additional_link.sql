-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

DO $$
BEGIN
    -- Check if the column already exists
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='additional_link'
    ) THEN
        ALTER TABLE events
        ADD COLUMN additional_link TEXT;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

ALTER TABLE events
DROP COLUMN IF EXISTS additional_link;
-- +goose StatementEnd
