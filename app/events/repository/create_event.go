package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) CreateEvent(ctx context.Context, event domain.Event) (uint, error) {
	err := repo.db.DB(ctx).Create(&event).Error
	if err != nil {
		logrus.Error("failed to create event")
		return 0, err
	}
	return event.ID, nil
}
