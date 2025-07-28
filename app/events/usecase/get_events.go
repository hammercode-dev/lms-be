package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) GetEvents(ctx context.Context, filter domain.EventFilter) (resp []domain.Event, pagination domain.Pagination, err error) {
	tData, datas, err := uc.repository.GetEvents(ctx, filter)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return
	}

	baseURL := config.GetConfig().BaseURL

	for i, data := range datas {
		datas[i].ImageEvent = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, data.ImageEvent)
	}

	return datas, domain.NewPagination(tData, filter.FilterPagination), err
}
