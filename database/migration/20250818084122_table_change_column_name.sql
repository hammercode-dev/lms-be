-- +goose Up
-- +goose StatementBegin
DO $$
BEGIN
    -- Check if github column exists and github_url doesn't
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='github'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='github_url'
    ) THEN
        ALTER TABLE "public"."users" RENAME COLUMN "github" TO "github_url";
    END IF;

    -- Check if linkedin column exists and linkedin_url doesn't
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='linkedin'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='linkedin_url'
    ) THEN
        ALTER TABLE "public"."users" RENAME COLUMN "linkedin" TO "linkedin_url";
    END IF;

    -- Check if personal_web column exists and personal_web_url doesn't
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='personal_web'
    ) AND NOT EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='personal_web_url'
    ) THEN
        ALTER TABLE "public"."users" RENAME COLUMN "personal_web" TO "personal_web_url";
    END IF;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
    -- Only rename back if the new columns exist and old ones don't
    IF EXISTS (
        SELECT 1 FROM information_schema.columns 
        WHERE table_name='users' AND column_name='github_url'
    ) THEN
        ALTER TABLE "public"."users" RENAME COLUMN "github_url" TO "github";
        ALTER TABLE "public"."users" RENAME COLUMN "linkedin_url" TO "linkedin";
        ALTER TABLE "public"."users" RENAME COLUMN "personal_web_url" TO "personal_web";
    END IF;
END $$;
-- +goose StatementEnd
