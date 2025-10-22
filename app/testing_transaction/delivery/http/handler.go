package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.TestingTransactionUsecase
}

// NewHandler - simple handler
func NewHandler(usecase domain.TestingTransactionUsecase) domain.TestingTransactionHandler {
	return &Handler{usecase: usecase}
}

// CreatePayment - endpoint untuk create payment
func (h *Handler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse request
	body, _ := io.ReadAll(r.Body)
	var req domain.CreatePaymentRequest
	if err := json.Unmarshal(body, &req); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Invalid request",
		}, w)
		return
	}

	// Call usecase
	resp, err := h.usecase.CreatePayment(ctx, req)
	if err != nil {
		logrus.WithError(err).Error("Failed to create payment")
		errMsg := err.Error()
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: "Failed to create payment: " + errMsg,
		}, w)
		return
	}

	// Success response
	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Payment created successfully",
		Data:    resp,
	}, w)
}

// GetPayment - endpoint untuk get payment by order_no
func (h *Handler) GetPayment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	orderNo := mux.Vars(r)["order_no"]

	// Call usecase
	transaction, err := h.usecase.GetPayment(ctx, orderNo)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    404,
			Message: "Payment not found",
		}, w)
		return
	}

	// Success response
	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Success",
		Data:    transaction,
	}, w)
}

// GetAllPayments - endpoint untuk get semua payments
func (h *Handler) GetAllPayments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Call usecase
	transactions, err := h.usecase.GetAllPayments(ctx)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: "Failed to get payments",
		}, w)
		return
	}

	// Success response
	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Success",
		Data:    transactions,
	}, w)
}

// XenditWebhook - endpoint untuk menerima webhook dari Xendit
// Xendit akan hit endpoint ini ketika payment berhasil/expired
func (h *Handler) XenditWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logrus.Info("Received Xendit webhook")

	// Parse webhook payload
	body, _ := io.ReadAll(r.Body)
	var payload map[string]interface{}
	if err := json.Unmarshal(body, &payload); err != nil {
		logrus.WithError(err).Error("Failed to parse webhook payload")
		w.WriteHeader(400)
		return
	}

	// Extract data dari webhook
	orderNo := payload["external_id"].(string)
	status := payload["status"].(string)        // PAID, EXPIRED
	paymentMethod := payload["payment_method"].(string) // BANK_TRANSFER, CREDIT_CARD, dll

	logrus.WithFields(logrus.Fields{
		"order_no":       orderNo,
		"status":         status,
		"payment_method": paymentMethod,
	}).Info("Processing webhook")

	// Convert Xendit status to our status
	ourStatus := "pending"
	if status == "PAID" {
		ourStatus = "paid"
	} else if status == "EXPIRED" {
		ourStatus = "expired"
	}

	// Call usecase
	if err := h.usecase.HandleWebhook(ctx, orderNo, ourStatus, paymentMethod); err != nil {
		logrus.WithError(err).Error("Failed to handle webhook")
		w.WriteHeader(500)
		return
	}

	// Return success to Xendit
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}
