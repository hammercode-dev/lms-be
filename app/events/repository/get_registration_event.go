package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetRegistrationEvent(ctx context.Context, orderNo string) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}

func (repo *repository) GetRegistrationEventByID(ctx context.Context, id uint) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Where("id = ?", id).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}
