package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

func (r *repository) CreateBlogPost(ctx context.Context, data domain.BlogPost) (domain.BlogPost, error) {
	if err := r.db.DB(ctx).Omit("updated_at").Create(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
