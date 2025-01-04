package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
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
// @Router /api/events/pay [post]
func (h Handler) PayEvent(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("failed to read body : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	var payload domain.EventPayPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		logrus.Error("failed to unmarshal : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}
	err = h.usecase.CreatePayEvent(r.Context(), payload)
	if err != nil {
		logrus.Error("failed to Create pay event : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    201,
		Message: "success",
		Data:    nil,
	}, w)
}
