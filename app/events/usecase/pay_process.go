package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) PayProcess(ctx context.Context, payload domain.PayProcessPayload) error {
	rEvent, err := uc.repository.GetRegistrationEvent(ctx, payload.OrderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if rEvent.ID == 0 {
		err = utils.NewNotFoundError(ctx, "registration order not found", errors.New("registration order not found"))
		return err
	}

	if rEvent.Status == "SUCCESS" {
		return nil
	}

	eventPay, err := uc.repository.GetEventPay(ctx, payload.OrderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if eventPay.ID == 0 {
		err = utils.NewNotFoundError(ctx, "event pay order not found", errors.New("event pay order not found"))
		return err
	}

	if eventPay.Status == "SUCCESS" {
		return nil
	}

	eventPay.Status = payload.Status
	rEvent.Status = payload.Status
	rEvent.UpToYou = payload.Note

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		err = uc.repository.UpdateEventPay(txCtx, eventPay)

		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		err = uc.repository.UpdateRegistrationEvent(txCtx, rEvent)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		return nil
	})

	return err
}
