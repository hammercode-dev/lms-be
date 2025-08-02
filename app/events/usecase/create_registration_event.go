package usecase

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/email"
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

	// Read HTML template for email
	htmlTmpl, err := os.ReadFile("./assets/event_status_registration_template.html")
	if err != nil {
		return domain.RegisterEventResponse{}, fmt.Errorf("failed to read file template: %w", err)
	}

	// Parse the HTML template
	tmpl, err := template.New("register_event_status").Parse(string(htmlTmpl))
	if err != nil {
		return domain.RegisterEventResponse{}, fmt.Errorf("failed to parse template: %w", err)
	}

	// Prepare data for the email template
	smtpConfig := email.SMTP{
		Email:    uc.cfg.SMTP_EMAIL,
		Password: uc.cfg.SMTP_PASSWORD,
		Host:     uc.cfg.SMTP_HOST,
		Port:     uc.cfg.SMTP_PORT,
	}

	// Prepare the receiver data
	emailPayload := email.NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		"Welcome to Our Platform - Event Registration Status",
		tmpl,
	)
	var formattedDate string
	if event.Date.Valid {
		formattedDate = event.Date.Time.Format("Monday, 02 January 2006")
	} else {
		formattedDate = "Date to be announced"
	}

	// Add the receiver's email and data to the payload
	if err := emailPayload.AddReceiver(
		ctx,
		email.Receiver{
			Email: payload.Email,
			Data: map[string]interface{}{
				"name":     payload.Name,
				"title":    event.Title,
				"price":    event.Price,
				"email":    payload.Email,
				"order_no": orderNo,
				"year":     time.Now().Format("2006"),
				"date":     formattedDate,
				"duration": event.Duration,
				"location": event.Location,
			},
		}); err != nil {
		return domain.RegisterEventResponse{}, fmt.Errorf("failed to add receiver: %w", err)
	}

	// is free event or not
	status := "SUCCESS"
	upToYou := "registration success"
	if event.Price != 0.0 {
		status = "PENDING"
		upToYou = "new register"
		emailPayload.SendEmail(ctx)
	} else {
		logrus.Info("free event, send email registration success")
		// Read HTML template for email
		htmlTmpl, err := os.ReadFile("./assets/event_status_registration_sucess_template.html")
		if err != nil {
			return domain.RegisterEventResponse{}, fmt.Errorf("failed to read file template: %w", err)
		}

		// Parse the HTML template
		tmpl, err := template.New("register_event_status").Parse(string(htmlTmpl))
		if err != nil {
			return domain.RegisterEventResponse{}, fmt.Errorf("failed to parse template: %w", err)
		}

		if err := emailPayload.ChangeTemplate(ctx, tmpl); err != nil {
			return domain.RegisterEventResponse{}, fmt.Errorf("failed to change template: %w", err)
		}
		emailPayload.SendEmail(ctx)
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
