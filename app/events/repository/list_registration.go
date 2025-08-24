package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) ListRegistration(ctx context.Context, filter domain.EventFilter, email string) (tData int, data []domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{}).
		Preload("Event").
		Preload("Event.Tags").
		Preload("Event.Speakers")

	if filter.Status != "" {
		db = db.Where("registration_events.status = ?", filter.Status)
	}

	if filter.StartDate.Valid {
		db = db.Where("registration_events.created_at >= ?", filter.StartDate)
	}

	if filter.EndDate.Valid {
		db = db.Where("registration_events.created_at <= ?", filter.EndDate)
	}

	if email != "" {
		db = db.Where("registration_events.email = ?", email)
	}

	var totalData int64
	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Find(&data).Error
	if err != nil {
		return
	}

	return int(totalData), data, err
}

