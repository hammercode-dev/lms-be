-- +goose Up
-- +goose StatementBegin

-- Buat ENUM blog_status hanya jika belum ada
DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1 FROM pg_type WHERE typname = 'blog_status'
    ) THEN
        CREATE TYPE blog_status AS ENUM ('draft', 'published', 'archived');
    END IF;
END$$;

-- Buat tabel authors jika belum ada
CREATE TABLE IF NOT EXISTS authors (
    user_id int8 PRIMARY KEY REFERENCES users(id),
    name VARCHAR(255) NOT NULL,
    avatar VARCHAR(512)
);

-- Buuat table sequence jika belum ada
CREATE SEQUENCE IF NOT EXISTS blog_posts_id_seq;
CREATE SEQUENCE IF NOT EXISTS blog_post_tags_id_seq;

-- Buat tabel blog_posts jika belum ada
CREATE TABLE IF NOT EXISTS blog_posts (
    id int8 PRIMARY KEY DEFAULT nextval('blog_posts_id_seq'),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    excerpt TEXT,
    author_id INTEGER REFERENCES authors(user_id),
    published_at TIMESTAMP default NULL,
    updated_at TIMESTAMP default NULL,
    category VARCHAR(100),
    status blog_status,
    slug VARCHAR(255) UNIQUE NOT NULL
);

-- Buat tabel blog_post_tags jika belum ada
CREATE TABLE IF NOT EXISTS blog_post_tags (
    id int8 PRIMARY KEY DEFAULT nextval('blog_post_tags_id_seq'),
    blog_post_id INTEGER REFERENCES blog_posts(id) ON DELETE CASCADE,
    tag VARCHAR(50)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS blog_post_tags;
DROP TABLE IF EXISTS blog_posts;

-- Hapus enum hanya jika tidak digunakan oleh kolom mana pun
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_type WHERE typname = 'blog_status') THEN
        -- Pastikan ENUM tidak digunakan oleh kolom lain sebelum DROP
        BEGIN
            DROP TYPE blog_status;
        EXCEPTION
            WHEN dependent_objects_still_exist THEN
                RAISE NOTICE 'blog_status is still in use, cannot drop.';
        END;
    END IF;
END$$;

DROP TABLE IF EXISTS authors;
-- +goose StatementEnd
