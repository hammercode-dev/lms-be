package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetEvents(ctx context.Context, filter domain.EventFilter) (tData int, dataDTO []domain.EventDTO, err error) {

	data := []domain.Event{}

	db := repo.db.DB(ctx).Model(&data)

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
		Preload("Tags").Preload("Speakers").Preload("Author").Find(&data).Error
	if err != nil {
		return
	}

	if err == nil {
		for _, d := range data {
			dataDTO = append(dataDTO, domain.EventDTO{
				ID:                   d.ID,
				Title:                d.Title,
				Description:          d.Description,
				Slug:                 d.Slug,
				Image:                d.Image,
				Date:                 d.Date,
				Type:                 d.Type,
				Location:             d.Location,
				Duration:             d.Duration,
				Capacity:             d.Capacity,
				Status:               d.Status,
				Tags:                 d.Tags,
				Speakers:             d.Speakers,
				SessionType:          d.SessionType,
				RegistrationLink:     d.RegistrationLink,
				Price:                d.Price,
				ReservationStartDate: d.ReservationStartDate,
				ReservationEndDate:   d.ReservationEndDate,
				Author:               d.Author.Username,
				CreatedAt:            d.CreatedAt,
				UpdatedAt:            d.UpdatedAt,
				DeletedAt:            d.DeletedAt,
			})
		}
	}

	return int(totalData), dataDTO, err
}
