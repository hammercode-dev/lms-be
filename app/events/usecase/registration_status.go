package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) RegistrationStatus(ctx context.Context, orderNo string) (resp domain.RegisterStatusResponse, err error) {
	rEvent, err := uc.repository.GetRegistrationEvent(ctx, orderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return resp, err
	}

	if rEvent.ID == 0 {
		err = utils.NewNotFoundError(ctx, "registration order not found", errors.New("registration order not found"))
		return resp, err
	}

	return domain.RegisterStatusResponse{
		OrderNo: orderNo,
		Status:  rEvent.Status,
	}, err
}
