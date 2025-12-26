package usecase

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// DeleteBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) DeleteBlogPost(ctx context.Context, id uint) error {
	if err := uc.repo.DeleteBlogPost(ctx, id); err != nil {
		logrus.Error("failed to delete blog post detail: ", err)
		return err
	}

	return nil
}
