package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) GetRegistrationEvent(ctx context.Context, orderNo string) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get registration event use generic conditions")
		return
	}

	return data, err
}
