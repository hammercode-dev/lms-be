package usecase

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"time"

	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/pkg/email"
	"github.com/hammer-code/lms-be/pkg/hash"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v4"
)

func (uc usecase) CreateRegistrationEvent(ctx context.Context, payload domain.RegisterEventPayload) (domain.RegisterEventResponse, error) {
	// get user data from context
	userData := ctx.Value(contextkey.UserKey).(domain.User)

	registrations, err := uc.repository.GetRegistrationEventUserByStatus(ctx, payload.EventID, uint(userData.ID), []string{constants.PENDING, constants.SUCCESS})
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return domain.RegisterEventResponse{}, err
	}

	if len(registrations) != 0 {
		err = utils.NewBadRequestError(ctx, "you have already registered", errors.New("you have already registered"))
		return domain.RegisterEventResponse{}, err
	}

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

	// check image proof payment
	dataImage := domain.Image{}
	if payload.ImageProofPayment != "" {
		dataImage, err = uc.imageRepository.GetImage(ctx, payload.ImageProofPayment)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return domain.RegisterEventResponse{}, err
		}

		if dataImage.IsUsed {
			err = utils.NewNotFoundError(ctx, "image not exists", errors.New("image not exists"))
			return domain.RegisterEventResponse{}, err
		}
	}

	// generate order number
	// format: TXE-<event_id>-<year><month><day><hash
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
			Email: userData.Email,
			Data: map[string]interface{}{
				"name":     userData.Username,
				"title":    event.Title,
				"price":    event.Price,
				"email":    userData.Email,
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
	status := constants.SUCCESS
	if event.Price != 0.0 {
		status = constants.PENDING
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
			OrderNo:           orderNo,
			UserID:            strconv.Itoa(userData.ID),
			EventID:           event.ID,
			Status:            status,
			ImageProofPayment: dataImage.FileName,
			PaymentDate:       null.NewTime(time.Now(), true),
		})

		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		if dataImage.FileName != "" {
			_, err = uc.repository.CreateEventPay(txCtx, domain.EventPay{
				RegistrationEventID: rId,
				EventID:             event.ID,
				OrderNO:             orderNo,
				ImageProofPayment:   dataImage.FileName,
				NetAmount:           payload.NetAmount,
			})

			if err != nil {
				err = utils.NewInternalServerError(ctx, err)
				return err
			}
		}

		err = uc.imageRepository.UpdateUseImage(txCtx, dataImage.ID)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}
		return nil
	})

	return domain.RegisterEventResponse{
		OrderNo: orderNo,
	}, err
}
