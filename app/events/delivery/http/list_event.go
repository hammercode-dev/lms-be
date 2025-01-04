package http

import (
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// GetEvents
// @Summary Get Events
// @Description This endpoint use to get events by filter
// @Tags Event
// @Accept json
// @Produce json
// @Param limit query string true "string"
// @Param page query string true "string"
// @Param orderBy query string false "string"
// @Param start_date query string false "string"
// @Param end_date query string false "string"
// @Param title query string false "string"
// @Param type query string false "string"
// @Param status query string false "string"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} []domain.Event
// @Router /api/events [get]
func (h Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
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

	data, pagination, err := h.usecase.GetEvents(r.Context(), domain.EventFilter{
		Title:            r.URL.Query().Get("title"),
		Type:             r.URL.Query().Get("type"),
		Status:           r.URL.Query().Get("status"),
		StartDate:        startDate,
		EndDate:          endDate,
		FilterPagination: flterPagination,
	})

	if err != nil {
		logrus.Error("failed to get event : ", err)
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
