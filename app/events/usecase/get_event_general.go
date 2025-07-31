package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (uc usecase) GetEventByID(ctx context.Context, id uint) (domain.Event, error) {
	resp, err := uc.repository.GetEvent(ctx, id)
	if err != nil {
		logrus.Error("failed to get event by id: ", err)
		return resp, err
	}
	baseURL := config.GetConfig().BaseURL

	resp.Image = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, resp.Image)
	return resp, err
}
