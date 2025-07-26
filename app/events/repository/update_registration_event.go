package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) UpdateRegistrationEvent(ctx context.Context, event domain.RegistrationEvent) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update registration event pay")
		return err
	}
	return nil
}
