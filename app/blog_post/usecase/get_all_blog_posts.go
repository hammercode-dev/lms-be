package usecase

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetAllBlogPosts implements domain.BlogPostUsecase.
func (uc *usecase) GetAllBlogPosts(ctx context.Context, pagination domain.FilterPagination) ([]domain.BlogPost, domain.Pagination, error) {
	// ensure default ordering lives in usecase layer
	if pagination.GetOrderBy() == "" {
		pagination.SetOrderBy("id DESC")
	}

	blogPosts, totalCount, err := uc.repo.GetAllBlogPosts(ctx, pagination)
	if err != nil {
		logrus.Error("failed to get all blog posts: ", err)
		return nil, domain.Pagination{}, err
	}

	// Batch fetch tags to avoid N+1 queries
	ids := make([]int, 0, len(blogPosts))
	for i := range blogPosts {
		ids = append(ids, blogPosts[i].Id)
	}

	if len(ids) > 0 {
		rawTags, tagErr := uc.repo.GetTagsByBlogPostIDs(ctx, ids)
		if tagErr != nil {
			logrus.Error("failed to batch get tags: ", tagErr)
		} else {
			// Transform []BlogPostTag to map for easier lookup
			tagsMap := make(map[int][]string)
			for _, t := range rawTags {
				tagsMap[t.BlogPostId] = append(tagsMap[t.BlogPostId], t.Tag)
			}
			// Attach tags to blog posts
			for i := range blogPosts {
				if tags, ok := tagsMap[blogPosts[i].Id]; ok {
					blogPosts[i].Tags = tags
				} else {
					blogPosts[i].Tags = nil
				}
			}
		}
	}

	return blogPosts, domain.NewPagination(totalCount, pagination), nil
}
