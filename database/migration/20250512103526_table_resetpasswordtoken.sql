-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE IF NOT EXISTS resetpasswordtoken_id_seq;

CREATE TABLE IF NOT EXISTS "public"."resetpasswordtoken" (
    "id" int8 NOT NULL DEFAULT nextval('resetpasswordtoken_id_seq'::regclass),
    "user_id" int8 NOT NULL,
    "token" varchar(255) NOT NULL,
    "expiry_date" timestamp NOT NULL,
    "is_used" boolean DEFAULT FALSE,
    "created_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT "resetpasswordtoken_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users"("id") ON DELETE CASCADE,
    PRIMARY KEY ("id")
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS "public"."resetpasswordtoken";
DROP SEQUENCE IF EXISTS resetpasswordtoken_id_seq;
-- +goose StatementEnd
