package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetEvent(ctx context.Context, eventID uint) (dataDTO domain.EventDTO, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	data := domain.Event{}

	err = db.Where("id = ?", eventID).Preload("Tags").Preload("Speakers").Preload("Author").First(&data).Error

	if err == nil {
		dataDTO = domain.EventDTO{
			ID:                   data.ID,
			Title:                data.Title,
			Description:          data.Description,
			Slug:                 data.Slug,
			Image:                data.Image,
			Date:                 data.Date,
			Type:                 data.Type,
			Location:             data.Location,
			Duration:             data.Duration,
			Capacity:             data.Capacity,
			Status:               data.Status,
			Tags:                 data.Tags,
			Speakers:             data.Speakers,
			SessionType:          data.SessionType,
			RegistrationLink:     data.RegistrationLink,
			Price:                data.Price,
			ReservationStartDate: data.ReservationStartDate,
			ReservationEndDate:   data.ReservationEndDate,
			Author:               data.Author.Username,
			CreatedAt:            data.CreatedAt,
			UpdatedAt:            data.UpdatedAt,
			DeletedAt:            data.DeletedAt,
		}
	}

	return dataDTO, err
}
