-- +goose Up
-- +goose StatementBegin

TRUNCATE TABLE "public"."users" RESTART IDENTITY CASCADE;

INSERT INTO "public"."users" (
    "username", "email", "password", "role", "fullname", 
    "date_of_birth", "gender", "phone_number", "address",
    "github", "linkedin", "personal_web", "created_at", "updated_at"
) VALUES 
 -- password : passowrd 
    ('admin', 'admin@example.com', '$2a$10$zzJJ6MKKBgJT0CfT7rjnWeCAfSRIG6VhBdoqSIWi1VjwBfsp6XcT.', 'admin', 'Admin User', 
     '1990-01-01', 'Male', '123456789', '123 Admin St',
     'github.com/admin', 'linkedin.com/in/admin', 'admin.com', CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
     
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM "public"."users" WHERE "username" IN ('admin', 'johndoe', 'janedoe');
-- +goose StatementEnd
