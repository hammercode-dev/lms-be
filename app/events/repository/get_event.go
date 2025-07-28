package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetEvent(ctx context.Context, eventID uint) (data domain.Event, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	err = db.Where("id = ?", eventID).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}