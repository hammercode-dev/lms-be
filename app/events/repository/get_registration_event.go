package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetRegistrationEvent(ctx context.Context, orderNo string) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Preload("User").Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}

func (repo *repository) GetRegistrationEventByID(ctx context.Context, id uint) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Preload("User").Where("id = ?", id).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}

func (repo *repository) GetRegistrationEventUserByStatus(ctx context.Context, eventID uint, userID string, statuses []string) (data []domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Where("event_id = ?", eventID).Where("user_id = ?", userID).Where("status IN (?)", statuses).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}
