package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (r repository) CreateRegistrationEvent(ctx context.Context, event domain.RegistrationEvent) (uint, error) {
	err := repo.db.DB(ctx).Create(&event).Error
	if err != nil {
		return 0, err
	}
	return event.ID, nil
}