package usecase

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// UpdateBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) UpdateBlogPost(ctx context.Context, data domain.BlogPost, id uint) error {
	err := uc.repo.UpdateBlogPost(ctx, data, id)
	if err != nil {
		logrus.Error("failed to update blog post: ", err)
		return err
	}
	return nil
}
