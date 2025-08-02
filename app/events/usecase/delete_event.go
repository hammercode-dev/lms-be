package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) DeleteEvent(ctx context.Context, id uint) error {
	err := uc.repository.DeleteEvent(ctx, id)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	return err
}