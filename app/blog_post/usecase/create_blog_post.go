package usecase

import (
	"errors"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

// CreateBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) CreateBlogPost(ctx context.Context, data domain.BlogPost, user domain.User) error {

	data.Author.UserId = user.ID
	data.Author.Name = user.Username
	data.UpdatedAt = nil
	data.PublishedAt = nil

	if data.Status == "published" {
		timeNow := time.Now()
		data.PublishedAt = &timeNow
	}

	if err := uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.repo.FindAuthorByUserID(txCtx, uint(user.ID))
		if errors.Is(err, gorm.ErrRecordNotFound) {

			if err := uc.repo.CreateAuthor(txCtx, data.Author); err != nil {
				return err
			}
		}
		data.AuthorID = data.Author.UserId

		data, err := uc.repo.CreateBlogPost(txCtx, data)
		if err != nil {
			return err
		}

		tags := make([]domain.BlogPostTag, 0, len(data.Tags))

		for _, tag := range data.Tags {
			tags = append(tags, domain.BlogPostTag{
				BlogPostId: data.Id,
				Tag:        tag,
			})
		}

		if err := uc.repo.CreateTags(txCtx, tags); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	return nil

}
