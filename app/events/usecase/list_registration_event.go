package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) ListRegistration(ctx context.Context, filter domain.EventFilter) (resp []domain.RegistrationEvent, pagination domain.Pagination, err error) {
	tData, datas, err := uc.repository.ListRegistration(ctx, filter)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return
	}

	baseURL := config.GetConfig().BaseURL

	for i, data := range datas {
		datas[i].ImageProofPayment = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, data.ImageProofPayment)
	}
	
	return datas, domain.NewPagination(tData, filter.FilterPagination), err
}
