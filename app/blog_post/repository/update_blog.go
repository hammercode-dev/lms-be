package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// UpdateBlogPost implements domain.BlogPostRepository.
func (r *repository) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	return r.db.StartTransaction(ctx, func(txCtx context.Context) error {
		if err := r.db.DB(txCtx).Model(&domain.BlogPost{}).
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
		if data.Author.Avatar != "" {
			if err := r.db.DB(txCtx).Model(&domain.Author{}).
				Where("user_id = ?", data.AuthorID).
				Updates(map[string]interface{}{
					"avatar": data.Author.Avatar,
				}).Error; err != nil {
				logrus.Error("failed to update author avatar: ", err)
				return err
			}
		}

		if len(data.Tags) > 0 {
			if err := r.db.DB(txCtx).Table("blog_post_tags").
				Where("blog_post_id = ?", id).
				Delete(nil).Error; err != nil {
				logrus.Error("failed to delete old tags: ", err)
				return err
			}

			for _, tag := range data.Tags {
				blogPostTag := struct {
					BlogPostId int    `gorm:"column:blog_post_id"`
					Tag        string `gorm:"column:tag"`
				}{
					BlogPostId: int(id),
					Tag:        tag,
				}

				if err := r.db.DB(txCtx).Table("blog_post_tags").
					Create(&blogPostTag).Error; err != nil {
					logrus.Error("failed to create blog post tag: ", err)
					return err
				}
			}
		}

		return nil
	})
}
