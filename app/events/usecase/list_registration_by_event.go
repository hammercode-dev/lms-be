package usecase

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) ListRegistrationByEvent(ctx context.Context, id uint, filterPagination domain.FilterPagination) (resp []domain.EventRegistrationDTO, pagination domain.Pagination, err error) {

	registrations, totalCount, err := uc.repository.ListRegistrationByEvent(ctx, id, filterPagination)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return resp, pagination, err
	}

	if len(registrations) == 0 {
		err = utils.NewNotFoundError(ctx, "registration event not found", errors.New("registration event not found"))
		return resp, pagination, err
	}

	for _, registration := range registrations {
		resp = append(resp, registration.ToRegistrationDTO())
	}

	pagination = domain.NewPagination(int(totalCount), filterPagination)

	return
}
