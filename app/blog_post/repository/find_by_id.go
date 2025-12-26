package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

// FindById implements domain.BlogPostRepository.
func (r *repository) FindById(ctx context.Context, id uint) (data domain.BlogPost, err error) {
	db := r.db.DB(ctx).Preload("Author").Model(&domain.BlogPost{}).Where("is_deleted = ?", false)
	err = db.First(&data, "id = ?", id).Error

	return data, err
}
