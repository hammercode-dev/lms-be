package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/utils"
	"gopkg.in/guregu/null.v4"
)

func (uc usecase) UpdateRegistrationStatus(ctx context.Context, id uint, payload domain.UpdateRegistrationStatusRequest) error {
	registrations, err := uc.repository.GetRegistrationEventByID(ctx, id)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if registrations.ID == 0 {
		err = utils.NewNotFoundError(ctx, "registration event not found", errors.New("registration event not found"))
		return err
	}

	userData := ctx.Value(contextkey.UserKey).(domain.User)

	registrations.Status = payload.Status
	registrations.UpdatedAt = null.TimeFrom(time.Now())
	registrations.UpdatedByUserID = userData.ID

	err = uc.repository.UpdateRegistrationEvent(ctx, registrations)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	return nil
}
