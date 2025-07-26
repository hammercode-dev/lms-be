package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) GetEventPay(ctx context.Context, orderNo string) (data domain.EventPay, err error) {
	db := repo.db.DB(ctx).Model(&domain.EventPay{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get event pay")
		return
	}

	return data, err
}