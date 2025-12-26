package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

func (r *repository) CreateTags(ctx context.Context, tag []domain.BlogPostTag) error {
	if err := r.db.DB(ctx).Table("blog_post_tags").Create(&tag).Error; err != nil {
		return err
	}
	return nil
}
