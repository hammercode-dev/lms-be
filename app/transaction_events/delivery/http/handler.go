package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

type Handler struct {
	usecase domain.TransactionEventUsecase
	cfg     config.Config
}

func NewHandler(usecase domain.TransactionEventUsecase, cfg config.Config) domain.TransactionEventHandler {
	return &Handler{
		usecase: usecase,
		cfg:     cfg,
	}
}

// CreateTransaction
// @Summary Create Payment Transaction
// @Description Create a payment transaction for event registration
// @Tags Transaction
// @Accept json
// @Produce json
// @Param request body domain.CreateTransactionPayload true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 201 {object} domain.HttpResponse
// @Router /api/v1/transactions [post]
func (h *Handler) CreateTransaction(w http.ResponseWriter, r *http.Request) {
	// Extract user from context
	userValue := r.Context().Value(contextkey.UserKey)
	if userValue == nil {
		ngelog.Error(r.Context(), "user not found in context", nil)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}, w)
		return
	}

	user, ok := userValue.(domain.User)
	if !ok {
		ngelog.Error(r.Context(), "failed to cast user from context", nil)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "Internal server error",
		}, w)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	var payload domain.CreateTransactionPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		ngelog.Error(r.Context(), "failed to unmarshal payload", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: err.Error(),
		}, w)
		return
	}

	data, err := h.usecase.CreateTransaction(r.Context(), user, payload)
	if err != nil {
		ngelog.Error(r.Context(), "failed to create transaction", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	message := "Please complete your payment to finalize your registration"
	if data.Amount == 0 {
		message = "Congratulations! You have been successfully registered for the event"
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: message,
		Data:    data,
	}, w)
}

// CheckPaymentStatus
// @Summary Check Payment Status
// @Description Check the status of a payment transaction
// @Tags Transaction
// @Accept json
// @Produce json
// @Param transaction_no path string true "Transaction Number"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/transactions/{transaction_no}/status [get]
func (h *Handler) CheckPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	transactionNo := vars["transaction_no"]

	if transactionNo == "" {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "transaction_no is required",
		}, w)
		return
	}

	data, err := h.usecase.CheckPaymentStatus(r.Context(), transactionNo)
	if err != nil {
		ngelog.Error(r.Context(), "failed to check payment status", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}, w)
}

// XenditWebhook
// @Summary Xendit Payment Webhook
// @Description Receive payment notification callback from Xendit
// @Tags Webhook
// @Accept json
// @Produce json
// @Param request body object true "Xendit Callback Payload"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /webhooks/xendit [post]
func (h *Handler) XenditWebhook(w http.ResponseWriter, r *http.Request) {
	// Verify X-CALLBACK-TOKEN
	callbackToken := r.Header.Get("X-CALLBACK-TOKEN")
	expectedToken := h.cfg.XENDIT_WEBHOOK_TOKEN

	if expectedToken != "" && callbackToken != expectedToken {
		ngelog.Error(r.Context(), "invalid webhook token", nil)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}, w)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read webhook body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}, w)
		return
	}

	// Log webhook for debugging
	ngelog.Info(r.Context(), "Received Xendit webhook: "+string(bodyBytes))

	// Try to parse as Invoice callback first
	var invoiceCallback domain.XenditInvoiceCallback
	if err := json.Unmarshal(bodyBytes, &invoiceCallback); err == nil && invoiceCallback.ID != "" {
		// Process as Invoice callback
		err = h.usecase.HandleXenditWebhook(r.Context(), invoiceCallback)
		if err != nil {
			ngelog.Error(r.Context(), "failed to process invoice webhook", err)
			resp := utils.CustomErrorResponse(err)
			utils.Response(resp, w)
			return
		}

		utils.Response(domain.HttpResponse{
			Code:    http.StatusOK,
			Message: "Webhook processed successfully",
		}, w)
		return
	}

	// Try to parse as Virtual Account callback
	var vaCallback domain.XenditVACallback
	if err := json.Unmarshal(bodyBytes, &vaCallback); err == nil && vaCallback.ExternalID != "" {
		// Process as VA callback
		err = h.usecase.HandleXenditWebhook(r.Context(), vaCallback)
		if err != nil {
			ngelog.Error(r.Context(), "failed to process VA webhook", err)
			resp := utils.CustomErrorResponse(err)
			utils.Response(resp, w)
			return
		}

		utils.Response(domain.HttpResponse{
			Code:    http.StatusOK,
			Message: "Webhook processed successfully",
		}, w)
		return
	}

	// Unknown webhook format
	ngelog.Error(r.Context(), "unknown webhook format", nil)
	utils.Response(domain.HttpResponse{
		Code:    http.StatusBadRequest,
		Message: "Unknown webhook format",
	}, w)
}

// GetOrderDetail
// @Summary Get Order Detail
// @Description Get detailed information about an order by order_no
// @Tags Transaction
// @Accept json
// @Produce json
// @Param order_no path string true "Order Number"
// @Failure 400 {object} domain.HttpResponse
// @Failure 404 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/orders/{order_no} [get]
func (h *Handler) GetOrderDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	orderNo := vars["order_no"]

	if orderNo == "" {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "order_no is required",
		}, w)
		return
	}

	data, err := h.usecase.GetOrderDetail(r.Context(), orderNo)
	if err != nil {
		ngelog.Error(r.Context(), "failed to get order detail", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}, w)
}