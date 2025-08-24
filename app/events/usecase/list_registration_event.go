package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) ListRegistration(ctx context.Context, filter domain.EventFilter) (resp []domain.RegistrationEvent, pagination domain.Pagination, err error) {
	userData := ctx.Value(contextkey.UserKey).(domain.User)

	tData, datas, err := uc.repository.ListRegistration(ctx, filter, userData.Email)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return
	}

	baseURL := config.GetConfig().BaseURL

	for i, data := range datas {
		datas[i].ImageProofPayment = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, data.ImageProofPayment)
		
		// Update event image URL
		if datas[i].Event.Image != "" {
			datas[i].Event.Image = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, datas[i].Event.Image)
		}
	}
	
	return datas, domain.NewPagination(tData, filter.FilterPagination), err
}
