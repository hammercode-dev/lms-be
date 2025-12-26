package usecase

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetAllBlogPosts implements domain.BlogPostUsecase.
func (uc *usecase) GetAllBlogPosts(ctx context.Context, pagination domain.FilterPagination) ([]domain.BlogPost, domain.Pagination, error) {
	blogPosts, totalCount, err := uc.repo.GetAllBlogPosts(ctx, pagination)
	if err != nil {
		logrus.Error("failed to get all blog posts: ", err)
		return nil, domain.Pagination{}, err
	}
	return blogPosts, domain.NewPagination(totalCount, pagination), nil
}
