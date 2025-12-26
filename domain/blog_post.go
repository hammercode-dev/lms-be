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
	CreateBlogPost(ctx context.Context, data BlogPost, user User) error
	UpdateBlogPost(ctx context.Context, data BlogPost, id uint) error
	DeleteBlogPost(ctx context.Context, id uint) error
	GetAllBlogPosts(ctx context.Context, pagination FilterPagination) ([]BlogPost, Pagination, error)
	GetDetailBlogPost(ctx context.Context, slug string, id uint) (BlogPost, error)
}

type BlogPostRepository interface {
	CreateBlogPost(ctx context.Context, data BlogPost) (BlogPost, error)
	UpdateBlogPost(ctx context.Context, data BlogPost, id uint) error
	DeleteBlogPost(ctx context.Context, id uint) error
	GetAllBlogPosts(ctx context.Context, pagination FilterPagination) ([]BlogPost, int, error)
	FindById(ctx context.Context, id uint) (BlogPost, error)
	FindBySlug(ctx context.Context, slug string) (BlogPost, error)
	GetTagsByBlogPostID(ctx context.Context, blogPostID uint) (tags []string, err error)
	FindAuthorByUserID(ctx context.Context, userID uint) (Author, error)
	CreateAuthor(ctx context.Context, data Author) error
	CreateTags(ctx context.Context, tag []BlogPostTag) error
}

type BlogPost struct {
	Id          int        `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Author      Author     `json:"author" gorm:"foreignKey:AuthorID;references:UserId"`
	AuthorID    int        `json:"author_id" gorm:"column:author_id"`
	Tags        []string   `json:"tags" gorm:"-"`
	Category    string     `json:"category"`
	Status      string     `json:"status" gorm:"type:enum('draft', 'published', 'archived')"`
	Slug        string     `json:"slug"`
	PublishedAt *time.Time `json:"published_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CreatedAt   *time.Time `json:"created_at"`
	IsDeleted   bool       `json:"-"`
}

type BlogPostTag struct {
	Id         int    `json:"id" gorm:"primaryKey"`
	BlogPostId int    `json:"blog_post_id"`
	Tag        string `json:"tag"`
}

type Author struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}
