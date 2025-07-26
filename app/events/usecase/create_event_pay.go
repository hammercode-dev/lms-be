package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (uc usecase) CreateEventPay(ctx context.Context, payload domain.EventPayPayload) error {
	rEvent, err := uc.repository.GetRegistrationEvent(ctx, payload.OrderNo)
	if err != nil {
		logrus.Error("failed to get event")
		return err
	}

	if rEvent.ID == 0 {
		return errors.New("registration order not found")
	}

	dataImage, err := uc.imageRepository.GetImage(ctx, payload.ImageProofPayment)
	if err != nil {
		logrus.Error("failed to create event", dataImage)
		return err
	}

	if dataImage.IsUsed {
		err = errors.New("image not exists")
		return err
	}

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		_, err := uc.repository.CreateEventPay(txCtx, domain.EventPay{
			RegistrationEventID: rEvent.ID,
			EventID:             rEvent.EventID,
			ImageProofPayment:   payload.ImageProofPayment,
			NetAmount:           payload.NetAmount,
			Status:              "PENDING",
			OrderNO:             rEvent.OrderNo,
		})

		if err != nil {
			logrus.Error("failed to create pay event")
			return err
		}

		err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
		if err != nil {
			logrus.Error("failed to update use image")
			return err
		}

		return nil
	})

	return err
}
