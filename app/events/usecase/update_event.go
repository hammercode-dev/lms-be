package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"gopkg.in/guregu/null.v4"
)

func (uc usecase) UpdateEvent(ctx context.Context, id uint, payload domain.UpdateEventPayload) error {
	if !constants.IsValidEventType(payload.Type) {
		return utils.NewBadRequestError(ctx, "Sorry, invalid event type", errors.New("event type is not valid"))
	}

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
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	return err
}
