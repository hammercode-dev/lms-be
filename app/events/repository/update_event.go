package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (r repository) UpdateEvent(ctx context.Context, event domain.Event) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		return err
	}

	return nil
}
