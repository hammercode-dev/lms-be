package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) DeleteEvent(ctx context.Context, id uint) error {
	err := repo.db.DB(ctx).Model(&domain.Event{}).Delete("id = ?", id).Error
	if err != nil {
		logrus.Error("failed to delete event")
		return err
	}
	return nil
}
