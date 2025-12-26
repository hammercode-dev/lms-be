package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

// GetTagsByBlogPostIDs returns raw blog post tags for given IDs.
func (r *repository) GetTagsByBlogPostIDs(ctx context.Context, blogPostIDs []int) ([]domain.BlogPostTag, error) {
	var tags []domain.BlogPostTag
	if len(blogPostIDs) == 0 {
		return tags, nil
	}
	if err := r.db.DB(ctx).Table("blog_post_tags").
		Where("blog_post_id IN (?)", blogPostIDs).
		Find(&tags).Error; err != nil {
		return nil, err
	}
	return tags, nil
}
