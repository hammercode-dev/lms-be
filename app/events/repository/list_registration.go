package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) ListRegistration(ctx context.Context, filter domain.EventFilter, email string) (tData int, data []domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{}).
		Preload("Event").
		Preload("Event.Tags").
		Preload("Event.Speakers").
		Preload("User").
		Preload("Transaction")

	// Filter by user email (only show registrations for the logged-in user)
	if email != "" {
		db = db.Joins(`
			JOIN users
			ON users.id = CAST(NULLIF(registration_events.user_id, '') AS INTEGER)
		`).Where("users.email = ?", email)
	}

	if filter.Status != "" {
		db = db.Where("registration_events.status = ?", filter.Status)
	}

	if filter.Type != "" {
		db = db.Joins("JOIN events ON events.id = registration_events.event_id").
			Where("events.type = ?", filter.Type)
	}

	if filter.StartDate.Valid {
		db = db.Where("registration_events.created_at >= ?", filter.StartDate)
	}

	if filter.EndDate.Valid {
		db = db.Where("registration_events.created_at <= ?", filter.EndDate)
	}

	var totalData int64
	// Count after all filters
	if err := db.Count(&totalData).Error; err != nil {
		return 0, nil, err
	}

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Find(&data).Error
	if err != nil {
		return
	}

	return int(totalData), data, nil
}
