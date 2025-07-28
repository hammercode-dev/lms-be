package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) ListRegistration(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	var totalData int64

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
		Offset(filter.FilterPagination.GetOffset()).Find(&data).Error
	if err != nil {
		return
	}

	return int(totalData), data, err
}