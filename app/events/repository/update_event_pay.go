package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (r repository) UpdateEventPay(ctx context.Context, event domain.EventPay) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		return err
	}
	return nil
}