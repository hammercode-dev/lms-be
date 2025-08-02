package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (uc usecase) GetEvents(ctx context.Context, filter domain.EventFilter) (resp []domain.Event, pagination domain.Pagination, err error) {
	tData, datas, err := uc.repository.GetEvents(ctx, filter)
	if err != nil {
		logrus.Error("failed to get event")
		return
	}

	baseURL := config.GetConfig().BaseURL

	for i, data := range datas {
		datas[i].Image = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, data.Image)
	}

	return datas, domain.NewPagination(tData, filter.FilterPagination), err
}
