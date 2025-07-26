package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) UpdateEventPay(ctx context.Context, event domain.EventPay) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update event pay")
		return err
	}
	return nil
}