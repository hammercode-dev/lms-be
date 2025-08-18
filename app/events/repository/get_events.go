package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetEvents(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.Event, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	var totalData int64

	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
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

	if filter.Title != "" {
		db = db.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Preload("Tags").Preload("Speakers").Find(&data).Error
	if err != nil {
		return
	}

	return int(totalData), data, err
}

