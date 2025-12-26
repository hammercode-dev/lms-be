package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

// GetAllBlogPosts implements domain.BlogPostRepository.
func (r *repository) GetAllBlogPosts(ctx context.Context, pagination domain.FilterPagination) ([]domain.BlogPost, int, error) {
	var data []domain.BlogPost
	var totalCount int64

	if err := r.db.DB(ctx).Model(&domain.BlogPost{}).Where("is_deleted = ?", false).Count(&totalCount).Error; err != nil {
		return nil, 0, err
	}

	query := r.db.DB(ctx).Preload("Author").Where("is_deleted = ?", false)

	if orderBy := pagination.GetOrderBy(); orderBy != "" {
		query = query.Order(orderBy)
	}

	err := query.Limit(pagination.GetLimit()).Offset(pagination.GetOffset()).Find(&data).Error
	if err != nil {
		return nil, 0, err
	}

	return data, int(totalCount), nil
}
