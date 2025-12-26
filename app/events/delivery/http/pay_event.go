package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

// PayEvent
// @Summary Pay Event
// @Description This endpoint pay the event
// @Tags Event
// @Accept json
// @Produce json
// @Param request body domain.EventPayPayload true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/events/pay [post]
func (h Handler) PayEvent(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	var payload domain.EventPayPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		ngelog.Error(r.Context(), "failed to unmarshal payload", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}
	err = h.usecase.CreateEventPay(r.Context(), payload)
	if err != nil {
		ngelog.Error(r.Context(), "failed to create pay event", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "success",
		Data:    nil,
	}, w)
}
