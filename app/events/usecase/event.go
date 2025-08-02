package usecase

import (
	"context"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

func (uc usecase) DeleteEvent(ctx context.Context, id uint) error {
	err := uc.repository.DeleteEvent(ctx, id)
	if err != nil {
		logrus.Error("failed to delete event by id: ", err)
		return err
	}

	return err
}

func (uc usecase) UpdateEvent(ctx context.Context, id uint, payload domain.UpdateEventPayload) error {
	err := uc.repository.UpdateEvent(ctx, domain.Event{
		ID:                   id,
		Title:                payload.Title,
		Description:          payload.Description,
		Slug:                 payload.Slug,
		Author:               payload.Author,
		Image:                payload.FileName,
		Date:                 payload.Date,
		Type:                 payload.Type,
		Location:             payload.Location,
		Duration:             payload.Duration,
		Capacity:             payload.Capacity,
		Status:               payload.Status,
		RegistrationLink:     payload.RegistrationLink,
		Price:                payload.Price,
		ReservationStartDate: payload.ReservationStartDate,
		ReservationEndDate:   payload.ReservationEndDate,
		UpdatedAt:            null.TimeFrom(time.Now()),
	})
	if err != nil {
		logrus.Error("failed to update event by id: ", err)
		return err
	}

	return err
}
