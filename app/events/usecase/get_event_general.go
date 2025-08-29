package usecase

import (
	"context"
	"fmt"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) GetEventByID(ctx context.Context, id uint) (resp domain.EventDTO, err error) {
	event, err := uc.repository.GetEvent(ctx, id)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return resp, err
	}
	baseURL := config.GetConfig().BaseURL

	event.Image = fmt.Sprintf("%s/api/v1/public/storage/images/%s", baseURL, event.Image)
	
	return event.ToDTO(), err
}
