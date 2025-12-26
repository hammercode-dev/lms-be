package repository

import (
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DeleteBlogPost implements domain.BlogPostRepository.
func (r *repository) DeleteBlogPost(ctx context.Context, id uint) error {
	db := r.db.DB(ctx).Model(&domain.BlogPost{})

	// Perform soft delete by updating is_deleted field
	result := db.Where("id = ?", id).Updates(map[string]interface{}{
		"is_deleted": true,
	})

	if result.Error != nil {
		logrus.Error("failed to soft delete blog post: ", result.Error)
		return result.Error
	}

	if result.RowsAffected == 0 {
		logrus.Warn("no blog post found to delete with id: ", id)
		return errors.New("blog post not found")
	}
	return nil
}
