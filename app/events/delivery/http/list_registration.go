package http

import (
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// ListRegistration
// @Summary Get List Register user to event
// @Description This endpoint use to get list of user that register on event
// @Tags Event
// @Accept json
// @Produce json
// @Param limit query string true "string"
// @Param page query string true "string"
// @Param event_id query string false "string"
// @Param start_date query string false "string"
// @Param end_date query string false "string"
// @Param status query string false "string"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} []domain.RegistrationEvent
// @Router /api/events/registrations [get]
func (h Handler) ListRegistration(w http.ResponseWriter, r *http.Request) {
	flterPagination, err := domain.GetPaginationFromCtx(r)
	if err != nil {
		logrus.Error("failed to get pagination : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	startDate, _ := utils.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := utils.ParseDate(r.URL.Query().Get("end_date"))

	data, pagination, err := h.usecase.ListRegistration(r.Context(), domain.EventFilter{
		Status:           r.URL.Query().Get("status"),
		StartDate:        startDate,
		EndDate:          endDate,
		FilterPagination: flterPagination,
	})

	if err != nil {
		logrus.Error("failed to get registration event : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:       200,
		Message:    "success",
		Data:       data,
		Pagination: pagination,
	}, w)
}
