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
        -- Add the new author_id column first
        ALTER TABLE events ADD COLUMN author_id INT;
        
        -- You'll need to populate author_id based on your business logic
        -- For example, if you have a users table:
        -- UPDATE events SET author_id = users.id FROM users WHERE events.author = users.name;
        
        -- Or set a default value for now:
        -- UPDATE events SET author_id = 1 WHERE author_id IS NULL;
        
        -- Drop the old author column after migration
        -- ALTER TABLE events DROP COLUMN author;
        
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
DO $$
BEGIN
    -- Add back the author column if it doesn't exist
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='events' AND column_name='author'
    ) THEN
        ALTER TABLE events ADD COLUMN author TEXT;
        
        -- You might want to populate it from a users table:
        -- UPDATE events SET author = users.name FROM users WHERE events.author_id = users.id;
    END IF;
    
    -- Drop the columns we added
    ALTER TABLE events DROP COLUMN IF EXISTS author_id;
    ALTER TABLE events DROP COLUMN IF EXISTS session_type;
END $$;
-- +goose StatementEnd