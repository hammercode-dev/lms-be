-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Kolom status
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'logout' AND column_name = 'status'
    ) THEN
        ALTER TABLE public.logout ADD COLUMN status smallint DEFAULT 1;
    END IF;

    -- Kolom user_id
    IF NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name = 'logout' AND column_name = 'user_id'
    ) THEN
        ALTER TABLE public.logout ADD COLUMN user_id int8;
    END IF;

    -- Constraint foreign key
    BEGIN
        ALTER TABLE public.logout ADD CONSTRAINT fk_logout_user 
        FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;
    EXCEPTION WHEN duplicate_object THEN NULL;
    END;

    -- Constraint CHECK status
    BEGIN
        ALTER TABLE public.logout ADD CONSTRAINT chk_logout_status 
        CHECK (status IN (0, 1));
    EXCEPTION WHEN duplicate_object THEN NULL;
    END;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE public.logout
DROP CONSTRAINT IF EXISTS fk_logout_user;
ALTER TABLE public.logout
DROP CONSTRAINT IF EXISTS chk_logout_status;

ALTER TABLE public.logout
DROP COLUMN IF EXISTS user_id,
DROP COLUMN IF EXISTS status;
-- +goose StatementEnd
