package usecase

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetDetailBlogPost implements domain.BlogPostUsecase.
func (uc *usecase) GetDetailBlogPost(ctx context.Context, slug string, id uint) (data domain.BlogPost, err error) {
	typeFind := "slug"
	if slug == "" {
		typeFind = "id"
	}

	if typeFind == "id" {
		data, err = uc.repo.FindById(ctx, id)
		if err != nil {
			logrus.Error("failed to find blog post by id: ", err)
			return domain.BlogPost{}, err
		}
	} else {
		data, err = uc.repo.FindBySlug(ctx, slug)
		if err != nil {
			logrus.Error("failed to find blog post by slug: ", err)
			return domain.BlogPost{}, err
		}
	}

	tags, err := uc.repo.GetTagsByBlogPostID(ctx, uint(data.Id))
	if err != nil {
		logrus.Error("failed to get tags for blog post ID ", data.Id, ": ", err)
		return domain.BlogPost{}, err
	}

	data.Tags = tags

	return data, nil
}
