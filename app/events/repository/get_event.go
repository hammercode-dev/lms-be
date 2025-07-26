package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) GetEvent(ctx context.Context, eventID uint) (data domain.Event, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	err = db.Where("id = ?", eventID).Find(&data).Error
	if err != nil {
		logrus.Error("repo.GetEvents: failed to get events use generic conditions")
		return
	}

	return data, err
}