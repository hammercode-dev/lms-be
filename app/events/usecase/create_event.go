package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) CreateEvent(ctx context.Context, payload domain.CreateEventPayload) error {
	if !constants.IsValidEventType(payload.Type) {
		return utils.NewBadRequestError(ctx, "Sorry, invalid event type", errors.New("event type is not valid"))
	}

	dataImage, err := uc.imageRepository.GetImage(ctx, payload.FileName)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if dataImage.IsUsed {
		err = utils.NewNotFoundError(ctx, "image not exists", errors.New("image not exists"))
		return err
	}

	userData := ctx.Value(contextkey.UserKey).(domain.User)

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		data := domain.Event{
			Title:                payload.Title,
			Description:          payload.Description,
			AuthorID:             userData.ID,
			Image:                dataImage.FileName,
			Date:                 payload.Date,
			Slug:                 payload.Slug,
			Type:                 payload.Type,
			Location:             payload.Location,
			Duration:             payload.Duration,
			Capacity:             payload.Capacity,
			RegistrationLink:     payload.RegistrationLink,
			ReservationStartDate: payload.ReservationStartDate,
			ReservationEndDate:   payload.ReservationEndDate,
			Price:                payload.Price,
			Status:               payload.Status,
			AdditionalLink:       payload.AdditionalLink,
			SessionType: 					payload.SessionType,
		}

		eventID, err := uc.repository.CreateEvent(txCtx, data)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		for _, tag := range payload.Tags {
			_, err = uc.repository.CreateEventTag(txCtx, domain.EventTag{
				EventID: eventID,
				Tag:     tag,
			})
			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}
		}

		for _, speaker := range payload.Speakers {
			_, err = uc.repository.CreateEventSpeaker(txCtx, domain.EventSpeaker{
				EventID: eventID,
				Name:    speaker,
			})
			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}
		}
		return nil
	})

	return err
}
