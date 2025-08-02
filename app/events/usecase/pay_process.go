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
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

func (uc usecase) PayProcess(ctx context.Context, payload domain.PayProcessPayload) error {
	rEvent, err := uc.repository.GetRegistrationEvent(ctx, payload.OrderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	logrus.Info("registration event: ", rEvent)
	if rEvent.ID == 0 {
		err = utils.NewNotFoundError(ctx, "registration order not found", errors.New("registration order not found"))
		return err
	}

	if rEvent.Status == "SUCCESS" {
		return nil
	}

	eventPay, err := uc.repository.GetEventPay(ctx, payload.OrderNo)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	if eventPay.ID == 0 {
		err = utils.NewNotFoundError(ctx, "event pay order not found", errors.New("event pay order not found"))
		return err
	}

	if eventPay.Status == "SUCCESS" {
		return nil
	}

	eventPay.Status = payload.Status
	rEvent.Status = payload.Status
	rEvent.UpToYou = payload.Note

	err = uc.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		err = uc.repository.UpdateEventPay(txCtx, eventPay)

		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		err = uc.repository.UpdateRegistrationEvent(txCtx, rEvent)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		// Get the complete event details
		event, err := uc.repository.GetEvent(txCtx, rEvent.EventID)
		if err != nil {
			logrus.Error("failed to get event details")
			return err
		}

		// Format the date in a readable format
		var formattedDate string
		if event.Date.Valid {
			formattedDate = event.Date.Time.Format("Monday, 02 January 2006")
		} else {
			formattedDate = "Date to be announced"
		}

		// Read HTML template for email
		var htmlTmpl []byte
		if payload.Status == "SUCCESS" {
			htmlTmpl, err = os.ReadFile("./assets/event_status_registration_sucess_template.html")
		} else {
			htmlTmpl, err = os.ReadFile("./assets/event_status_registration_template.html")
		}
		if err != nil {
			return fmt.Errorf("failed to read file template: %w", err)
		}

		// Parse the HTML template
		tmpl, err := template.New("register_event_status").Parse(string(htmlTmpl))
		if err != nil {
			return fmt.Errorf("failed to parse template: %w", err)
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
			"Update: Your Event Registration Status",
			tmpl,
		)

		// Add the receiver's email and data to the payload
		if err := emailPayload.AddReceiver(
			ctx,
			email.Receiver{
				Email: rEvent.Email,
				Data: map[string]interface{}{
					"name":     rEvent.Name,
					"title":    event.Title,
					"price":    event.Price,
					"email":    rEvent.Email,
					"order_no": rEvent.OrderNo,
					"status":   payload.Status,
					"note":     payload.Note,
					"year":     time.Now().Format("2006"),
					"date":     formattedDate,
					"duration": event.Duration,
					"location": event.Location,
				},
			}); err != nil {
			return fmt.Errorf("failed to add receiver: %w", err)
		}

		// Send the email
		emailPayload.SendEmail(ctx)
		return nil
	})

	return err
}
