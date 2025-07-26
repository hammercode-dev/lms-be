package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) UpdateEvent(ctx context.Context, event domain.Event) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update event")
		return err
	}

	return nil
}
