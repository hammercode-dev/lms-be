package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) CreateEventTag(ctx context.Context, data domain.EventTag) (uint, error) {
	err := repo.db.DB(ctx).Create(&data).Error
	if err != nil {
		logrus.Error("failed to create event tag")
		return 0, err
	}
	return data.ID, nil
}