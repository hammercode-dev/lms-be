package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/hash"
	"github.com/hammer-code/lms-be/utils"
)

func (uc usecase) CreateRegistrationEvent(ctx context.Context, payload domain.RegisterEventPayload) (domain.RegisterEventResponse, error) {
	event, err := uc.repository.GetEvent(ctx, payload.EventID)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return domain.RegisterEventResponse{}, err
	}

	if event.ID == 0 {
		err = utils.NewNotFoundError(ctx, "event not found", errors.New("event not found"))
		return domain.RegisterEventResponse{}, err
	}

	tNow := time.Now()

	if !event.ReservationStartDate.Valid {
		err = utils.NewBadRequestError(ctx, "event is not start to booking", errors.New("event is not start to booking"))
		return domain.RegisterEventResponse{}, err
	}

	if event.ReservationEndDate.Valid {
		if tNow.After(event.ReservationEndDate.Time) {
			err = utils.NewBadRequestError(ctx, "priode booking has ended", errors.New("priode booking has ended"))
			return domain.RegisterEventResponse{}, err
		}
	}

	hash := hash.GenerateHash(time.Now().Format("2006-01-02 15:04:05"))

	orderNo := fmt.Sprintf("TXE-%d-%s%s%s%s", event.ID, time.Now().Format("06"), time.Now().Format("01"), time.Now().Format("02"), hash[0:4])

	// is free event or not
	status := "SUCCESS"
	upToYou := "registration success"
	if event.Price != 0.0 {
		status = "PENDING"
		upToYou = "new register"
	}

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		rId, err := uc.repository.CreateRegistrationEvent(txCtx, domain.RegistrationEvent{
			OrderNo:     orderNo,
			EventID:     event.ID,
			Name:        payload.Name,
			Email:       payload.Email,
			PhoneNumber: payload.PhoneNumber,
			Status:      status,
			UpToYou:     upToYou,
		})

		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		if payload.ImageProofPayment != "" {
			dataImage, err := uc.imageRepository.GetImage(ctx, payload.ImageProofPayment)
			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}

			if dataImage.IsUsed {
				err = utils.NewNotFoundError(ctx, "image not exists", errors.New("image not exists"))
				return err
			}

			_, err = uc.repository.CreateEventPay(txCtx, domain.EventPay{
				RegistrationEventID: rId,
				EventID:             event.ID,
				ImageProofPayment:   payload.ImageProofPayment,
				NetAmount:           payload.NetAmount,
			})

			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}

			err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}
		}
		return nil
	})

	return domain.RegisterEventResponse{
		OrderNo: orderNo,
	}, err
}
