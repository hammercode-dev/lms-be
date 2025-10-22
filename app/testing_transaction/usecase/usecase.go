package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/xendit"
	"github.com/sirupsen/logrus"
)

type usecase struct {
	repo         domain.TestingTransactionRepository
	xenditClient *xendit.Client
}

// NewUsecase - simple usecase
func NewUsecase(repo domain.TestingTransactionRepository, xenditClient *xendit.Client) domain.TestingTransactionUsecase {
	return &usecase{
		repo:         repo,
		xenditClient: xenditClient,
	}
}

// CreatePayment - buat payment baru
func (uc *usecase) CreatePayment(ctx context.Context, req domain.CreatePaymentRequest) (domain.CreatePaymentResponse, error) {
	// 1. Generate Order Number
	orderNo := fmt.Sprintf("ORDER-%d", time.Now().Unix())

	logrus.WithFields(logrus.Fields{
		"order_no": orderNo,
		"amount":   req.Amount,
		"email":    req.CustomerEmail,
	}).Info("Creating payment")

	// 2. Panggil Xendit untuk buat invoice
	invoiceURL, err := uc.xenditClient.CreateInvoice(
		ctx,
		orderNo,
		req.CustomerEmail,
		req.Amount,
		fmt.Sprintf("Payment for %s", req.CustomerName),
	)
	if err != nil {
		// Safe error logging - handle nil pointer from Xendit SDK
		errMsg := "unknown error"
		if err != nil {
			errMsg = fmt.Sprintf("%v", err)
		}
		logrus.WithField("error", errMsg).Error("Failed to create Xendit invoice")
		return domain.CreatePaymentResponse{}, fmt.Errorf("failed to create invoice: %s", errMsg)
	}

	logrus.WithField("invoice_url", invoiceURL).Info("Xendit invoice created")

	// 3. Simpan ke database
	transaction := domain.TestingTransaction{
		OrderNo:       orderNo,
		CustomerName:  req.CustomerName,
		CustomerEmail: req.CustomerEmail,
		Amount:        req.Amount,
		Status:        "pending",
		InvoiceURL:    invoiceURL,
	}

	if err := uc.repo.Create(ctx, transaction); err != nil {
		logrus.WithError(err).Error("Failed to save transaction")
		return domain.CreatePaymentResponse{}, fmt.Errorf("failed to save transaction: %w", err)
	}

	logrus.WithField("order_no", orderNo).Info("Transaction saved successfully")

	// 4. Return response
	return domain.CreatePaymentResponse{
		OrderNo:    orderNo,
		InvoiceURL: invoiceURL,
		Amount:     req.Amount,
		Status:     "pending",
	}, nil
}

// GetPayment - ambil data payment
func (uc *usecase) GetPayment(ctx context.Context, orderNo string) (domain.TestingTransaction, error) {
	return uc.repo.GetByOrderNo(ctx, orderNo)
}

// GetAllPayments - ambil semua payment
func (uc *usecase) GetAllPayments(ctx context.Context) ([]domain.TestingTransaction, error) {
	return uc.repo.GetAll(ctx)
}

// HandleWebhook - handle webhook dari Xendit (dipanggil ketika user sudah bayar)
func (uc *usecase) HandleWebhook(ctx context.Context, orderNo string, status string, paymentMethod string) error {
	logrus.WithFields(logrus.Fields{
		"order_no":       orderNo,
		"status":         status,
		"payment_method": paymentMethod,
	}).Info("Handling Xendit webhook")

	// Get transaction
	transaction, err := uc.repo.GetByOrderNo(ctx, orderNo)
	if err != nil {
		logrus.WithError(err).Error("Transaction not found")
		return fmt.Errorf("transaction not found: %w", err)
	}

	// Update status
	transaction.Status = status
	transaction.PaymentMethod = paymentMethod

	if err := uc.repo.Update(ctx, transaction); err != nil {
		logrus.WithError(err).Error("Failed to update transaction")
		return fmt.Errorf("failed to update transaction: %w", err)
	}

	logrus.WithField("order_no", orderNo).Info("Transaction updated successfully")
	return nil
}
