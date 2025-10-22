package domain

import (
	"context"
	"net/http"
	"time"
)

// TestingTransaction - simple model untuk testing payment gateway
type TestingTransaction struct {
	ID            uint      `json:"id" gorm:"primarykey"`
	OrderNo       string    `json:"order_no"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	Amount        float64   `json:"amount"`
	Status        string    `json:"status"` // pending, paid, expired
	InvoiceURL    string    `json:"invoice_url"`
	PaymentMethod string    `json:"payment_method"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (TestingTransaction) TableName() string {
	return "testing_transaction"
}

// Request untuk create payment
type CreatePaymentRequest struct {
	CustomerName  string  `json:"customer_name" validate:"required"`
	CustomerEmail string  `json:"customer_email" validate:"required,email"`
	Amount        float64 `json:"amount" validate:"required,gt=0"`
}

// Response setelah create payment
type CreatePaymentResponse struct {
	OrderNo    string  `json:"order_no"`
	InvoiceURL string  `json:"invoice_url"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

// Repository interface - simple CRUD
type TestingTransactionRepository interface {
	Create(ctx context.Context, data TestingTransaction) error
	GetByOrderNo(ctx context.Context, orderNo string) (TestingTransaction, error)
	GetAll(ctx context.Context) ([]TestingTransaction, error)
	Update(ctx context.Context, data TestingTransaction) error
}

// Usecase interface
type TestingTransactionUsecase interface {
	CreatePayment(ctx context.Context, req CreatePaymentRequest) (CreatePaymentResponse, error)
	GetPayment(ctx context.Context, orderNo string) (TestingTransaction, error)
	GetAllPayments(ctx context.Context) ([]TestingTransaction, error)
	HandleWebhook(ctx context.Context, orderNo string, status string, paymentMethod string) error
}

// Handler interface
type TestingTransactionHandler interface {
	CreatePayment(w http.ResponseWriter, r *http.Request)
	GetPayment(w http.ResponseWriter, r *http.Request)
	GetAllPayments(w http.ResponseWriter, r *http.Request)
	XenditWebhook(w http.ResponseWriter, r *http.Request)
}
