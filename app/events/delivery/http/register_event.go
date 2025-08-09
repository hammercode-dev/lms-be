package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

// RegisterEvent
// @Summary Register Event
// @Description This endpoint use to register event
// @Tags Event
// @Accept json
// @Produce json
// @Param request body domain.RegisterEventPayload true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/events/registrations [post]
func (h Handler) RegisterEvent(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	token := utils.ExtractBearerToken(r)

	if err != nil {
		ngelog.Error(r.Context(), "failed to verify token", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusUnauthorized,
			Message: "Unauthorized",
		}, w)
		return
	}

	var payload domain.RegisterEventPayload
	if err := json.Unmarshal(bodyBytes, &payload); err != nil {
		ngelog.Error(r.Context(), "failed to unmarshal payload", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}
	data, err := h.usecase.CreateRegistrationEvent(r.Context(), payload, *token)
	if err != nil {
		ngelog.Error(r.Context(), "failed to create registration event", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "success",
		Data:    data,
	}, w)
}
