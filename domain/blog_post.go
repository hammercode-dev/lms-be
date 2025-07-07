package domain

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
)

type BlogPostHandler interface {
	CreateBlogPost(w http.ResponseWriter, r *http.Request)
	UpdateBlogPost(w http.ResponseWriter, r *http.Request)
	DeleteBlogPost(w http.ResponseWriter, r *http.Request)
	GetAllBlogPosts(w http.ResponseWriter, r *http.Request)
	GetDetailBlogPost(w http.ResponseWriter, r *http.Request)
}

type BlogPostUsecase interface {
	CreateBlogPost(ctx context.Context, data BlogPost, token string) (BlogPost, error)
	UpdateBlogPost(ctx context.Context, data BlogPost, id uint) error
	DeleteBlogPost(ctx context.Context, id uint) error
	GetAllBlogPosts(ctx context.Context) ([]BlogPost, error)
	GetDetailBlogPost(ctx context.Context, slug string, id uint) (BlogPost, error)
}

type BlogPostRepository interface {
	CreateBlogPost(ctx context.Context, data BlogPost) (BlogPost, error)
	UpdateBlogPost(ctx context.Context, data BlogPost, id uint) error
	DeleteBlogPost(ctx context.Context, id uint) error
	GetAllBlogPosts(ctx context.Context) ([]BlogPost, error)
	GetDetailBlogPost(ctx context.Context, slug, typeFind string, id uint) (BlogPost, error)
}

type BlogPost struct {
	Id          int       `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	Author      Author    `json:"author" gorm:"foreignKey:AuthorID;references:UserId"`
	AuthorID    int       `json:"author_id" gorm:"column:author_id"`
	PublishedAt time.Time `json:"published_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Tags        []string  `json:"tags" gorm:"-"`
	Category    string    `json:"category"`
	Status      string    `json:"status" gorm:"type:enum('draft', 'published', 'archived')"`
	Slug        string    `json:"slug"`
}

type Author struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
