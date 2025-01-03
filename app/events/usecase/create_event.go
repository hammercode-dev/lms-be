package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (uc usecase) CreateEvent(ctx context.Context, payload domain.CreateEventPayload) error {
	dataImage, err := uc.imageRepository.GetImage(ctx, payload.FileName)
	if err != nil {
		logrus.Error("failed to create event", dataImage)
		return err
	}

	if dataImage.IsUsed {
		err = errors.New("image not exists")
		return err
	}

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		data := domain.Event{
			Title:            payload.Title,
			Description:      payload.Description,
			Author:           payload.Author,
			ImageEvent:       dataImage.FileName,
			DateEvent:        payload.DateEvent,
			Type:             payload.Type,
			Location:         payload.Location,
			Duration:         payload.Duration,
			Capacity:         payload.Capacity,
			RegistrationLink: payload.RegistrationLink,
			BookingStart:     payload.BookingStart,
			BookingEnd:       payload.BookingEnd,
			Price:            payload.Price,
			Status:           payload.Status,
		}

		eventID, err := uc.repository.CreateEvent(txCtx, data)
		if err != nil {
			logrus.Error("failed to create event", data)
			return err
		}

		err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
		if err != nil {
			logrus.Error("failed to update use image", data)
			return err
		}

		for _, tag := range payload.Tags {
			_, err = uc.repository.CreateEventTag(txCtx, domain.EventTag{
				EventID: eventID,
				Tag:     tag,
			})
			if err != nil {
				logrus.Error("failed to create event tag", data)
				return err
			}
		}

		for _, speaker := range payload.Speakers {
			_, err = uc.repository.CreateEventSpeaker(txCtx, domain.EventSpeaker{
				EventID: eventID,
				Name:    speaker,
			})
			if err != nil {
				logrus.Error("failed to create event speaker", data)
				return err
			}
		}
		return nil
	})

	return err
}
