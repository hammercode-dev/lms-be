-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Check if author column exists and author_id doesn't
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='author'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='author_id'
    ) THEN
        -- Rename and change type to INT if needed
        ALTER TABLE events RENAME COLUMN author TO author_id;
        ALTER TABLE events ALTER COLUMN author_id TYPE INT USING (author_id::integer);
    ELSIF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='author_id'
    ) THEN
        -- If neither column exists, add author_id as INT
        ALTER TABLE events ADD COLUMN author_id INT;
    END IF;

    -- Check if session_type field exists, add it if it doesn't
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='session_type'
    ) THEN
        ALTER TABLE events ADD COLUMN session_type TEXT;
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE events RENAME COLUMN author_id TO author;
ALTER TABLE events DROP COLUMN IF EXISTS session_type;
-- +goose StatementEnd
