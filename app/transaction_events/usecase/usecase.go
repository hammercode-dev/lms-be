package usecase

import (
	"context"
	"errors"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/email"
	"github.com/hammer-code/lms-be/pkg/hash"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/pkg/xendit"
	"github.com/hammer-code/lms-be/utils"
	"gopkg.in/guregu/null.v4"
	"gorm.io/gorm"
)

type usecase struct {
	transactionRepo domain.TransactionEventRepository
	eventRepo       domain.EventRepository
	xenditClient    *xendit.Client
	cfg             config.Config
}

func NewUsecase(
	transactionRepo domain.TransactionEventRepository,
	eventRepo domain.EventRepository,
	xenditClient *xendit.Client,
	cfg config.Config,
) domain.TransactionEventUsecase {
	return &usecase{
		transactionRepo: transactionRepo,
		eventRepo:       eventRepo,
		xenditClient:    xenditClient,
		cfg:             cfg,
	}
}

func (u *usecase) CreateTransaction(ctx context.Context, user domain.User, payload domain.CreateTransactionPayload) (domain.CreateTransactionResponse, error) {
	// 1. Get event details
	event, err := u.eventRepo.GetEvent(ctx, payload.EventID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.CreateTransactionResponse{}, utils.NewNotFoundError(ctx, "Event not found", errors.New("event not found"))
		}
		return domain.CreateTransactionResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 2. Check if user already registered for this event
	userID := fmt.Sprintf("%d", user.ID)
	registrations, err := u.eventRepo.GetRegistrationEventUserByStatus(
		ctx,
		payload.EventID,
		userID,
		[]string{constants.PENDING, constants.SUCCESS},
	)
	if err != nil {
		return domain.CreateTransactionResponse{}, utils.NewInternalServerError(ctx, err)
	}

	if len(registrations) > 0 {
		return domain.CreateTransactionResponse{}, utils.NewBadRequestError(ctx, "You have already registered for this event", errors.New("duplicate registration"))
	}

	// 3. Generate order_no for registration
	hashStr := hash.GenerateHash(time.Now().Format("2006-01-02 15:04:05"))
	orderNo := fmt.Sprintf("TXE-%d-%s%s%s%s",
		event.ID,
		time.Now().Format("06"),
		time.Now().Format("01"),
		time.Now().Format("02"),
		hashStr[0:4],
	)

	// 4. Determine status based on event price
	status := constants.PENDING
	if event.Price == 0 {
		status = constants.SUCCESS // Free event, langsung success
	}

	// 5. Create registration_events
	registrationID, err := u.eventRepo.CreateRegistrationEvent(ctx, domain.RegistrationEvent{
		OrderNo:     orderNo,
		UserID:      userID,
		EventID:     payload.EventID,
		Status:      status,
		PaymentDate: null.NewTime(time.Now(), event.Price == 0), // set payment_date jika free event
	})
	if err != nil {
		return domain.CreateTransactionResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 6. If free event, send success email and return response
	if event.Price == 0 {
		// Send success registration email in background
		go u.sendRegistrationEmail(user, event, orderNo, true)

		return domain.CreateTransactionResponse{
			TransactionNo: "",
			OrderNo:       orderNo,
			Amount:        0,
			PaymentURL:    "",
			Status:        status,
		}, nil
	}

	// 7. For paid events, generate transaction_no
	transactionNo := fmt.Sprintf("TRX-%s-%d", time.Now().Format("20060102150405"), registrationID)

	// 8. Create invoice in Xendit with redirect URLs
	successRedirect := u.cfg.XENDIT_SUCCESS_REDIRECT
	if successRedirect != "" {
		if successRedirect[len(successRedirect)-1] != '/' {
			successRedirect += "/"
		}
		successRedirect += orderNo
	}

	// Use orderNo as external_id (first parameter) so Xendit can use it in redirect URL
	invoiceURL, invoiceID, err := u.xenditClient.CreateInvoice(
		ctx,
		orderNo, // Use orderNo as external_id for Xendit
		user.Email,
		event.Price,
		fmt.Sprintf("Payment for %s - %s", event.Title, orderNo),
		successRedirect,
		u.cfg.XENDIT_FAILURE_REDIRECT,
	)
	if err != nil {
		return domain.CreateTransactionResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 9. Save transaction to database
	transaction := domain.TransactionEvent{
		TransactionNo:  transactionNo,
		RegistrationID: registrationID,
		Amount:         event.Price,
		Status:         "pending",
		InvoiceURL:     &invoiceURL,
		InvoiceID:      &invoiceID,
		ExternalID:     &orderNo, // Store orderNo as external_id for webhook lookup
		CreatedAt:      time.Now(),
	}

	saved, err := u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return domain.CreateTransactionResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// Send pending payment email in background
	go u.sendRegistrationEmail(user, event, orderNo, false)

	return domain.CreateTransactionResponse{
		TransactionNo: saved.TransactionNo,
		OrderNo:       orderNo,
		Amount:        saved.Amount,
		PaymentURL:    *saved.InvoiceURL,
		Status:        saved.Status,
	}, nil
}

func (u *usecase) CheckPaymentStatus(ctx context.Context, transactionNo string) (domain.CheckPaymentStatusResponse, error) {
	// 1. Get transaction from database
	transaction, err := u.transactionRepo.GetByTransactionNo(ctx, transactionNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.CheckPaymentStatusResponse{}, utils.NewNotFoundError(ctx, "Transaction not found", errors.New("transaction not found"))
		}
		return domain.CheckPaymentStatusResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 2. If already paid, just return from database
	if transaction.Status == "paid" {
		return domain.CheckPaymentStatusResponse{
			TransactionNo: transaction.TransactionNo,
			Status:        transaction.Status,
			PaidAt:        transaction.PaidAt.Ptr(),
			PaymentMethod: transaction.PaymentMethod,
		}, nil
	}

	// 3. Check status from Xendit API
	if transaction.InvoiceID == nil || *transaction.InvoiceID == "" {
		return domain.CheckPaymentStatusResponse{}, utils.NewBadRequestError(ctx, "Transaction has no invoice ID", errors.New("invoice_id is empty"))
	}

	status, paidAtStr, paymentMethod, err := u.xenditClient.GetInvoiceStatus(ctx, *transaction.InvoiceID)
	if err != nil {
		return domain.CheckPaymentStatusResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 4. If status changed, update transaction_events
	if status != transaction.Status {
		transaction.Status = status

		// Parse paidAt string to time.Time
		if paidAtStr != nil {
			paidAtTime, err := time.Parse("2006-01-02 15:04:05", *paidAtStr)
			if err == nil {
				transaction.PaidAt = null.NewTime(paidAtTime, true)
			}
		}

		if paymentMethod != nil {
			transaction.PaymentMethod = paymentMethod
		}

		transaction.UpdatedAt = null.NewTime(time.Now(), true)

		// Update transaction
		err = u.transactionRepo.Update(ctx, transaction)
		if err != nil {
			return domain.CheckPaymentStatusResponse{}, utils.NewInternalServerError(ctx, err)
		}
	}

	// 5. Update registration_events status to match transaction_events status
	// This runs regardless of whether transaction status changed, to ensure sync
	registration, err := u.eventRepo.GetRegistrationEventByID(ctx, transaction.RegistrationID)
	if err != nil {
		return domain.CheckPaymentStatusResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// Determine what registration status should be based on transaction status
	var newRegistrationStatus string
	switch transaction.Status {
	case "paid":
		// Payment successful
		newRegistrationStatus = constants.SUCCESS
	case "expired":
		// Payment expired
		newRegistrationStatus = constants.EXPIRED
	case "failed":
		// Payment failed
		newRegistrationStatus = constants.CANCELLED
	default:
		// For "pending" or other statuses, keep as PENDING
		newRegistrationStatus = constants.PENDING
	}

	// Only update if registration status is different
	if registration.Status != newRegistrationStatus {
		registration.Status = newRegistrationStatus

		// Update payment_date if paid
		if transaction.Status == "paid" {
			registration.PaymentDate = transaction.PaidAt
		}

		registration.UpdatedAt = null.NewTime(time.Now(), true)

		err = u.eventRepo.UpdateRegistrationEvent(ctx, registration)
		if err != nil {
			return domain.CheckPaymentStatusResponse{}, utils.NewInternalServerError(ctx, err)
		}
	}

	return domain.CheckPaymentStatusResponse{
		TransactionNo: transaction.TransactionNo,
		Status:        transaction.Status,
		PaidAt:        transaction.PaidAt.Ptr(),
		PaymentMethod: transaction.PaymentMethod,
	}, nil
}

// HandleXenditWebhook processes webhook callback from Xendit
func (u *usecase) HandleXenditWebhook(ctx context.Context, callback interface{}) error {
	var externalID string
	var amount float64
	var status string
	var paidAtStr string
	var paymentChannel string

	// Handle different callback types (Invoice or Virtual Account)
	switch v := callback.(type) {
	case domain.XenditInvoiceCallback:
		externalID = v.ExternalID
		amount = v.Amount
		status = v.Status
		paidAtStr = v.PaidAt
		paymentChannel = v.PaymentChannel

	case domain.XenditVACallback:
		externalID = v.ExternalID
		amount = v.Amount
		status = "PAID" // VA callback only sent when paid
		paidAtStr = v.TransactionTimestamp
		paymentChannel = fmt.Sprintf("VA_%s", v.BankCode)

	default:
		return utils.NewBadRequestError(ctx, "Invalid callback type", errors.New("unknown callback type"))
	}

	// Only process PAID or SETTLED status
	if status != "PAID" && status != "SETTLED" {
		ngelog.Info(ctx, fmt.Sprintf("Webhook received but status is %s, skipping", status))
		return nil
	}

	// Get transaction by external_id
	transaction, err := u.transactionRepo.GetByExternalID(ctx, externalID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.NewNotFoundError(ctx, "Transaction not found", errors.New("transaction not found"))
		}
		return utils.NewInternalServerError(ctx, err)
	}

	// Check if already paid (idempotency)
	if transaction.Status == "paid" {
		ngelog.Info(ctx, "Transaction already paid, skipping update")
		return nil
	}

	// Verify amount matches
	if transaction.Amount != amount {
		return utils.NewBadRequestError(ctx, "Amount mismatch", fmt.Errorf("expected %.2f, got %.2f", transaction.Amount, amount))
	}

	// Parse paid_at timestamp
	paidAt, err := time.Parse(time.RFC3339, paidAtStr)
	if err != nil {
		// Try alternative format
		paidAt, err = time.Parse("2006-01-02T15:04:05.000Z", paidAtStr)
		if err != nil {
			ngelog.Error(ctx, "Failed to parse paid_at", err)
			paidAt = time.Now()
		}
	}

	// Update transaction status
	transaction.Status = "paid"
	transaction.PaidAt = null.NewTime(paidAt, true)
	transaction.PaymentMethod = &paymentChannel
	transaction.UpdatedAt = null.NewTime(time.Now(), true)

	err = u.transactionRepo.Update(ctx, transaction)
	if err != nil {
		return utils.NewInternalServerError(ctx, err)
	}

	// Update registration status
	registration, err := u.eventRepo.GetRegistrationEventByID(ctx, transaction.RegistrationID)
	if err != nil {
		return utils.NewInternalServerError(ctx, err)
	}

	registration.Status = constants.SUCCESS
	registration.PaymentDate = null.NewTime(paidAt, true)
	registration.UpdatedAt = null.NewTime(time.Now(), true)

	err = u.eventRepo.UpdateRegistrationEvent(ctx, registration)
	if err != nil {
		return utils.NewInternalServerError(ctx, err)
	}

	// Get user and event details for email
	event, err := u.eventRepo.GetEvent(ctx, registration.EventID)
	if err != nil {
		ngelog.Error(ctx, "Failed to get event for email", err)
		return nil // Don't fail webhook processing if email fails
	}

	// Send success email in background
	go u.sendRegistrationEmail(registration.User, event, registration.OrderNo, true)

	ngelog.Info(ctx, fmt.Sprintf("Webhook processed successfully for transaction %s", transaction.TransactionNo))
	return nil
}

func (u *usecase) GetOrderDetail(ctx context.Context, orderNo string) (domain.GetOrderDetailResponse, error) {
	// 1. Get registration by order_no
	registration, err := u.eventRepo.GetRegistrationEvent(ctx, orderNo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return domain.GetOrderDetailResponse{}, utils.NewNotFoundError(ctx, "Order not found", errors.New("order not found"))
		}
		return domain.GetOrderDetailResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 2. Get event details
	event, err := u.eventRepo.GetEvent(ctx, registration.EventID)
	if err != nil {
		return domain.GetOrderDetailResponse{}, utils.NewInternalServerError(ctx, err)
	}

	// 3. Get user details
	// User should already be preloaded in registration
	if registration.User.ID == 0 {
		return domain.GetOrderDetailResponse{}, utils.NewInternalServerError(ctx, errors.New("user not found in registration"))
	}

	// 4. Get transaction details (if exists)
	var transactionNo string
	transaction, err := u.transactionRepo.GetByRegistrationID(ctx, registration.ID)
	if err == nil {
		transactionNo = transaction.TransactionNo
	}
	// If error, it means no transaction (free event), transactionNo stays empty

	// 5. Build response
	var eventDate *time.Time
	if event.Date.Valid {
		eventDate = &event.Date.Time
	}

	var paymentDate *time.Time
	if registration.PaymentDate.Valid {
		paymentDate = &registration.PaymentDate.Time
	}

	response := domain.GetOrderDetailResponse{
		OrderNo:       registration.OrderNo,
		TransactionNo: transactionNo,
		PaymentDate:   paymentDate,
		Status:        registration.Status,
		EventDetail: domain.OrderEventDetail{
			Title:       event.Title,
			Date:        eventDate,
			Type:        string(event.Type),
			Location:    event.Location,
			Duration:    event.Duration,
			Price:       event.Price,
			SessionType: event.SessionType,
		},
		UserDetail: domain.OrderUserDetail{
			Fullname:    registration.User.Fullname,
			Email:       registration.User.Email,
			PhoneNumber: registration.User.PhoneNumber,
		},
	}

	return response, nil
}

// sendRegistrationEmail sends email notification for event registration
// isFree=true: send success registration email
// isFree=false: send pending payment email
func (u *usecase) sendRegistrationEmail(user domain.User, event domain.Event, orderNo string, isFree bool) {
	// Create new context for background operation
	ctx := context.Background()

	// Select template based on event type
	var templatePath string
	var subject string

	if isFree {
		templatePath = "./assets/event_status_registration_sucess_template.html"
		subject = "Registration Confirmed - " + event.Title
	} else {
		templatePath = "./assets/event_status_registration_template.html"
		subject = "Complete Your Payment - " + event.Title
	}

	// Read HTML template
	htmlTmpl, err := os.ReadFile(templatePath)
	if err != nil {
		ngelog.Error(ctx, "Failed to read email template", err)
		return
	}

	// Parse template
	tmpl, err := template.New("event_registration").Parse(string(htmlTmpl))
	if err != nil {
		ngelog.Error(ctx, "Failed to parse email template", err)
		return
	}

	// Prepare SMTP config
	smtpConfig := email.SMTP{
		Email:    u.cfg.SMTP_EMAIL,
		Password: u.cfg.SMTP_PASSWORD,
		Host:     u.cfg.SMTP_HOST,
		Port:     u.cfg.SMTP_PORT,
	}

	// Create email payload
	emailPayload := email.NewSendEmail(
		ctx,
		smtpConfig,
		"MIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8",
		subject,
		tmpl,
	)

	// Format event date
	var formattedDate string
	if event.Date.Valid {
		formattedDate = event.Date.Time.Format("Monday, 02 January 2006")
	} else {
		formattedDate = "Date to be announced"
	}

	// Add receiver with data
	err = emailPayload.AddReceiver(ctx, email.Receiver{
		Email: user.Email,
		Data: map[string]interface{}{
			"name":     user.Username,
			"title":    event.Title,
			"price":    event.Price,
			"email":    user.Email,
			"order_no": orderNo,
			"year":     time.Now().Format("2006"),
			"date":     formattedDate,
			"duration": event.Duration,
			"location": event.Location,
		},
	})

	if err != nil {
		ngelog.Error(ctx, "Failed to add email receiver", err)
		return
	}

	// Send email
	emailPayload.SendEmail(ctx)
}
