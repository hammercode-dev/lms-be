-- +goose Up
-- +goose StatementBegin
ALTER TABLE public.logout
ADD COLUMN status smallint DEFAULT 1,
ADD COLUMN user_id int8;

ALTER TABLE public.logout
ADD CONSTRAINT fk_logout_user FOREIGN KEY (user_id) REFERENCES public.users(id) ON DELETE SET NULL;

ALTER TABLE public.logout
ADD CONSTRAINT chk_logout_status CHECK (status IN (0, 1));
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
