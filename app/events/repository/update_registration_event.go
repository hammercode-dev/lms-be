package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (r repository) UpdateRegistrationEvent(ctx context.Context, event domain.RegistrationEvent) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		return err
	}
	return nil
}
