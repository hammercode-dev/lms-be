package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) CreateEventPay(ctx context.Context, payload domain.EventPayPayload) error {
	rEvent, err := uc.repository.GetRegistrationEvent(ctx, payload.OrderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if rEvent.ID == 0 {
		err = utils.NewNotFoundError(ctx, "registration order not found", errors.New("registration order not found"))
		return err
	}

	dataImage, err := uc.imageRepository.GetImage(ctx, payload.ImageProofPayment)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if dataImage.IsUsed {
		err = utils.NewNotFoundError(ctx, "image not exists", errors.New("image not exists"))
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
			err = utils.NewInternalServerError(txCtx, err)
			return err
		}

		err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
		if err != nil {
			err = utils.NewInternalServerError(txCtx, err)
			return err
		}

		return nil
	})

	return err
}
