package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

func (r *repository) FindAuthorByUserID(ctx context.Context, userID uint) (data domain.Author, err error) {
	if err := r.db.DB(ctx).Where("user_id = ?", userID).First(&data).Error; err != nil {
		return data, err
	}
	return data, nil
}
