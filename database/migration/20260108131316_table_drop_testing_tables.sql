-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."testing_transaction";
DROP TABLE IF EXISTS "public"."testing";
DROP SEQUENCE IF EXISTS testing_id_seq;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Recreate if needed for rollback (optional)
-- +goose StatementEnd