package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (r repository) DeleteEvent(ctx context.Context, id uint) error {
	err := repo.db.DB(ctx).Model(&domain.Event{}).Delete("id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}
