package repository

import (
	"github.com/hammer-code/lms-be/domain"
	"golang.org/x/net/context"
)

func (r *repository) CreateAuthor(ctx context.Context, data domain.Author) error {
	if err := r.db.DB(ctx).Create(&data).Error; err != nil {
		return err
	}
	return nil
}
