package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

// FindBySlug implements domain.BlogPostRepository.
func (r *repository) FindBySlug(ctx context.Context, slug string) (data domain.BlogPost, err error) {
	db := r.db.DB(ctx).Preload("Author").Model(&domain.BlogPost{}).Where("is_deleted = ?", false)
	err = db.First(&data, "slug = ?", slug).Error

	return data, err
}
