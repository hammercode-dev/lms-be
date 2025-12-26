package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

// DeleteBlogPost implements domain.BlogPostRepository.
func (r *repository) DeleteBlogPost(ctx context.Context, id uint) error {
	return r.db.DB(ctx).Model(&domain.BlogPost{}).Where("id = ?", id).Updates(map[string]interface{}{"is_deleted": true}).Error
}
