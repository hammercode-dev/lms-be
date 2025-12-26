package usecase

import (
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// UpdateBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	now := time.Now()
	data.UpdatedAt = &now
	if data.Status == "published" && data.PublishedAt == nil {
		data.PublishedAt = &now
	}

	if err := uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		if err := uc.repo.UpdateBlogPost(txCtx, data, id); err != nil {
			logrus.Error("failed to update blog post: ", err)
			return err
		}

		if data.Author.Avatar != "" {
			if err := uc.repo.UpdateAuthorAvatar(txCtx, uint(data.AuthorID), data.Author.Avatar); err != nil {
				logrus.Error("failed to update author avatar: ", err)
				return err
			}
		}

		if len(data.Tags) > 0 {
			if err := uc.repo.DeleteTagsByBlogPostID(txCtx, id); err != nil {
				logrus.Error("failed to delete old tags: ", err)
				return err
			}

			tags := make([]domain.BlogPostTag, 0, len(data.Tags))
			for _, tag := range data.Tags {
				tags = append(tags, domain.BlogPostTag{
					BlogPostId: int(id),
					Tag:        tag,
				})
			}

			if err := uc.repo.CreateTags(txCtx, tags); err != nil {
				logrus.Error("failed to create blog post tags: ", err)
				return err
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
