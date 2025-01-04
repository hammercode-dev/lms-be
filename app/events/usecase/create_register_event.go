package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/hash"
	"github.com/sirupsen/logrus"
)

func (uc usecase) CreateRegisterEvent(ctx context.Context, payload domain.RegisterEventPayload) (domain.RegisterEventResponse, error) {
	event, err := uc.repository.GetEvent(ctx, payload.EventID)
	if err != nil {
		logrus.Error("failed to get event")
		return domain.RegisterEventResponse{}, err
	}

	if event.ID == 0 {
		return domain.RegisterEventResponse{}, errors.New("event not found")
	}

	tNow := time.Now()

	if !event.ReservationStartDate.Valid {
		return domain.RegisterEventResponse{}, errors.New("event is not start to booking")
	}

	if event.ReservationEndDate.Valid {
		if tNow.After(event.ReservationEndDate.Time) {
			return domain.RegisterEventResponse{}, errors.New("priode booking has ended")
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
		rId, err := uc.repository.CreateRegisterEvent(txCtx, domain.RegistrationEvent{
			OrderNo:     orderNo,
			EventID:     event.ID,
			Name:        payload.Name,
			Email:       payload.Email,
			PhoneNumber: payload.PhoneNumber,
			Status:      status,
			UpToYou:     upToYou,
		})

		if err != nil {
			logrus.Error("failed to get event")
			return err
		}

		if payload.ImageProofPayment != "" {
			dataImage, err := uc.imageRepository.GetImage(ctx, payload.ImageProofPayment)
			if err != nil {
				logrus.Error("failed to create event", dataImage)
				return err
			}

			if dataImage.IsUsed {
				err = errors.New("image not exists")
				return err
			}

			_, err = uc.repository.CreatePayEvent(txCtx, domain.EventPay{
				RegistrationEventID: rId,
				EventID:             event.ID,
				ImageProofPayment:   payload.ImageProofPayment,
				NetAmount:           payload.NetAmount,
			})

			if err != nil {
				logrus.Error("failed to create pay event")
				return err
			}

			err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
			if err != nil {
				logrus.Error("failed to update use image")
				return err
			}
		}
		return nil
	})

	return domain.RegisterEventResponse{
		OrderNo: orderNo,
	}, err
}
