package domain

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type TransactionEvent struct {
    ID             uint      `json:"id" gorm:"primarykey"`
    TransactionNo  string    `json:"transaction_no" gorm:"uniqueIndex"`
    RegistrationID uint      `json:"registration_id"`
    Amount         float64   `json:"amount"`
    Status         string    `json:"status"` // pending, paid, expired
    InvoiceID      *string   `json:"invoice_id"`
    InvoiceURL     *string   `json:"invoice_url"`
    ExternalID     *string   `json:"external_id" gorm:"uniqueIndex"`
    PaymentMethod  *string   `json:"payment_method"`
    PaidAt         null.Time `json:"paid_at"`
    CreatedAt      time.Time `json:"created_at"`
    UpdatedAt      null.Time `json:"updated_at"`
    Registration   RegistrationEvent `json:"registration" gorm:"foreignKey:RegistrationID"`
}

func (TransactionEvent) TableName() string {
    return "transaction_events"
}

// Request & Response
type CreateTransactionPayload struct {
    EventID uint `json:"event_id"`
}

type CreateTransactionResponse struct {
    TransactionNo string  `json:"transaction_no"`
    OrderNo       string  `json:"order_no"`
    Amount        float64 `json:"amount"`
    PaymentURL    string  `json:"payment_url"`
    Status        string  `json:"status"`
}

type CheckPaymentStatusResponse struct {
    TransactionNo string  `json:"transaction_no"`
    Status        string  `json:"status"`
    PaidAt        *time.Time `json:"paid_at"`
    PaymentMethod *string `json:"payment_method"`
}

type GetOrderDetailResponse struct {
    OrderNo       string              `json:"order_no"`
    TransactionNo string              `json:"transaction_no"`
    PaymentDate   *time.Time          `json:"payment_date"`
    Status        string              `json:"status"`
    EventDetail   OrderEventDetail    `json:"event_detail"`
    UserDetail    OrderUserDetail     `json:"user_detail"`
}

type OrderEventDetail struct {
    Title       string    `json:"title"`
    Date        *time.Time `json:"date"`
    Type        string    `json:"type"`
    Location    string    `json:"location"`
    Duration    string    `json:"duration"`
    Price       float64   `json:"price"`
    SessionType string    `json:"session_type"`
}

type OrderUserDetail struct {
    Fullname    string `json:"fullname"`
    Email       string `json:"email"`
    PhoneNumber string `json:"phone_number"`
}

// Xendit Webhook Callbacks
type XenditInvoiceCallback struct {
    ID                     string  `json:"id"`
    ExternalID             string  `json:"external_id"`
    UserID                 string  `json:"user_id"`
    Status                 string  `json:"status"`
    MerchantName           string  `json:"merchant_name"`
    Amount                 float64 `json:"amount"`
    PayerEmail             string  `json:"payer_email"`
    Description            string  `json:"description"`
    PaidAt                 string  `json:"paid_at"`
    Updated                string  `json:"updated"`
    Created                string  `json:"created"`
    Currency               string  `json:"currency"`
    PaymentChannel         string  `json:"payment_channel"`
    PaymentDestination     string  `json:"payment_destination"`
}

type XenditVACallback struct {
    ID                          string  `json:"id"`
    PaymentID                   string  `json:"payment_id"`
    CallbackVirtualAccountID    string  `json:"callback_virtual_account_id"`
    ExternalID                  string  `json:"external_id"`
    AccountNumber               string  `json:"account_number"`
    BankCode                    string  `json:"bank_code"`
    Amount                      float64 `json:"amount"`
    TransactionTimestamp        string  `json:"transaction_timestamp"`
    MerchantCode                string  `json:"merchant_code"`
}

// Repository Interface
type TransactionEventRepository interface {
    Create(ctx context.Context, data TransactionEvent) (TransactionEvent, error)
    GetByTransactionNo(ctx context.Context, transactionNo string) (TransactionEvent, error)
    GetByRegistrationID(ctx context.Context, registrationID uint) (TransactionEvent, error)
    GetByExternalID(ctx context.Context, externalID string) (TransactionEvent, error)
    Update(ctx context.Context, data TransactionEvent) error
}

// Usecase Interface
type TransactionEventUsecase interface {
    CreateTransaction(ctx context.Context, user User, payload CreateTransactionPayload) (CreateTransactionResponse, error)
    CheckPaymentStatus(ctx context.Context, transactionNo string) (CheckPaymentStatusResponse, error)
    HandleXenditWebhook(ctx context.Context, callback interface{}) error
    GetOrderDetail(ctx context.Context, orderNo string) (GetOrderDetailResponse, error)
}

// Handler Interface
type TransactionEventHandler interface {
    CreateTransaction(w http.ResponseWriter, r *http.Request)
    CheckPaymentStatus(w http.ResponseWriter, r *http.Request)
    XenditWebhook(w http.ResponseWriter, r *http.Request)
    GetOrderDetail(w http.ResponseWriter, r *http.Request)
}