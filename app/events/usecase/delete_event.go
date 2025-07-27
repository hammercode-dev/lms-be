package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
)

func (uc usecase) DeleteEvent(ctx context.Context, id uint) error {
	err := uc.repository.DeleteEvent(ctx, id)
	if err != nil {
		logrus.Error("failed to delete event by id: ", err)
		return err
	}

	return err
}