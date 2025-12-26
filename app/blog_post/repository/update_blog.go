package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// UpdateBlogPost implements domain.BlogPostRepository.
func (r *repository) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	if err := r.db.DB(ctx).Model(&domain.BlogPost{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"title":        data.Title,
			"content":      data.Content,
			"excerpt":      data.Excerpt,
			"published_at": data.PublishedAt,
			"updated_at":   data.UpdatedAt,
			"category":     data.Category,
			"status":       data.Status,
		}).Error; err != nil {
		logrus.Error("failed to update blog post: ", err)
		return err
	}
	return nil
}
