package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

// GetAllBlogPosts implements domain.BlogPostRepository.
func (r *repository) GetAllBlogPosts(ctx context.Context, pagination domain.FilterPagination) ([]domain.BlogPost, int, error) {
	var data []domain.BlogPost
	var totalCount int64

	if err := r.db.DB(ctx).Model(&domain.BlogPost{}).Where("is_deleted = ?", false).Count(&totalCount).Error; err != nil {
		logrus.Error("failed to count blog posts: ", err)
		return nil, 0, err
	}

	offset := pagination.GetOffset()
	limit := pagination.GetLimit()
	orderBy := pagination.GetOrderBy()

	query := r.db.DB(ctx).Preload("Author").Where("is_deleted = ?", false)

	if orderBy != "" {
		query = query.Order(orderBy)
	} else {
		query = query.Order("id DESC")
	}

	err := query.Limit(limit).Offset(offset).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get all blog posts: ", err)
		return nil, 0, err
	}

	for i := range data {
		var tags []string
		if err := r.db.DB(ctx).Table("blog_post_tags").
			Select("tag").
			Where("blog_post_id = ?", data[i].Id).
			Pluck("tag", &tags).Error; err != nil {
			logrus.Error("failed to get tags for blog post ID ", data[i].Id, ": ", err)
		} else {
			data[i].Tags = tags
		}
	}

	return data, int(totalCount), nil
}
