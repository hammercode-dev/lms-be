package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetEventPay(ctx context.Context, orderNo string) (data domain.EventPay, err error) {
	db := repo.db.DB(ctx).Model(&domain.EventPay{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		return
	}

	return data, err
}