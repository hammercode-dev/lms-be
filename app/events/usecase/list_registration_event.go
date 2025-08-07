package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) ListRegistration(ctx context.Context, filter domain.EventFilter, token string) (resp []domain.RegistrationEvent, pagination domain.Pagination, err error) {

	userData, err := uc.jwt.VerifyToken(token)
	if err != nil {
		return nil, domain.Pagination{}, utils.NewUnauthorizedError(ctx, "unauthorized", err)
	}

	tData, datas, err := uc.repository.ListRegistration(ctx, filter, userData.Email)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return
	}

	return datas, domain.NewPagination(tData, filter.FilterPagination), err
}
