package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) ListEventPay(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.EventPay, err error) {
	db := repo.db.DB(ctx).Model(&domain.EventPay{})

	var totalData int64

	if filter.ID != 0 {
		db = db.Where("event_id = ?", filter.ID)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	if filter.StartDate.Valid {
		db = db.Where("start_date > ?", filter.StartDate)
	}

	if filter.StartDate.Valid {
		db = db.Where("end_date < ?", filter.EndDate)
	}

	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Preload("RegistrationEvent").Find(&data).Error
	if err != nil {
		logrus.Error("repo.GetEvents: failed to get event pays use generic conditions")
		return
	}

	return int(totalData), data, err
}