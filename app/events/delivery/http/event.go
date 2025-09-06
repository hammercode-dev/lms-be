package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

// @Summary Create Event
// @Description This endpoint create event
// @Tags Event
// @Accept json
// @Produce json
// @Param request body domain.CreateEventPayload true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/admin/events [post]
// @Security BearerAuth
func (h Handler) CreateEvent(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	var payload domain.CreateEventPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		ngelog.Error(r.Context(), "failed to unmarshal", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	err = h.usecase.CreateEvent(r.Context(), payload)
	if err != nil {
		ngelog.Error(r.Context(), "failed to create event", err)
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

// @Summary Update Event
// @Description This endpoint update event
// @Tags Event
// @Accept json
// @Produce json
// @Param request body domain.UpdateEventPayload true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/admin/events/{id} [patch]
// @Security BearerAuth
func (h Handler) UpdateEvent(w http.ResponseWriter, r *http.Request) {
	idS := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		ngelog.Error(r.Context(), "failed to convert string to uint", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ngelog.Error(r.Context(), "failed to read body", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	var payload domain.UpdateEventPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		ngelog.Error(r.Context(), "failed to unmarshal", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	err = h.usecase.UpdateEvent(r.Context(), uint(id), payload)
	if err != nil {
		ngelog.Error(r.Context(), "failed to update event", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
	}, w)
}

// @Summary Get Detail Event by ID
// @Description This endpoint use to get event by id
// @Tags Event
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.Event
// @Router /api/v1/events/{id} [get]
func (h Handler) GetEventByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		ngelog.Error(r.Context(), "failed to convert string to uint", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	data, err := h.usecase.GetEventByID(r.Context(), uint(id))
	if err != nil {
		ngelog.Error(r.Context(), "failed to get event by id", err)
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

// @Summary Get Detail Event by ID
// @Description This endpoint use to get detail event by id
// @Tags Event
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.Event
// @Router /api/v1/admin/events/{id} [get]
func (h Handler) GetDetail(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	id, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		ngelog.Error(r.Context(), "failed to convert string to uint", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	data, err := h.usecase.GetEventByIDAdmin(r.Context(), uint(id))
	if err != nil {
		ngelog.Error(r.Context(), "failed to get event by id", err)
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

// @Summary Delete Event By Id
// @Description This endpoint use to delete event by id
// @Tags Event
// @Accept json
// @Produce json
// @Param id path string true "id"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/v1/admin/events/{id} [delete]
// @Security BearerAuth
func (h Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	value, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		ngelog.Error(r.Context(), "failed to convert string to uint", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	err = h.usecase.DeleteEvent(r.Context(), uint(value))
	if err != nil {
		ngelog.Error(r.Context(), "failed to delete event", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "success",
		Data:    nil,
	}, w)
}

// @Summary Get List Available Event
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
// @Router /api/v1/events [get]
func (h Handler) List(w http.ResponseWriter, r *http.Request) {
	flterPagination, err := domain.GetPaginationFromCtx(r)
	if err != nil {
		ngelog.Error(r.Context(), "failed to get pagination", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	startDate, _ := utils.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := utils.ParseDate(r.URL.Query().Get("end_date"))

	data, pagination, err := h.usecase.GetEvents(r.Context(), domain.EventFilter{
		Title:            r.URL.Query().Get("title"),
		Type:             constants.EventType(r.URL.Query().Get("type")),
		Status:           r.URL.Query().Get("status"),
		StartDate:        startDate,
		EndDate:          endDate,
		FilterPagination: flterPagination,
	})

	if err != nil {
		ngelog.Error(r.Context(), "failed to get events", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:       http.StatusOK,
		Message:    "success",
		Data:       data,
		Pagination: &pagination,
	}, w)
}

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
// @Router /api/v1/admin/events [get]
func (h Handler) GetEvents(w http.ResponseWriter, r *http.Request) {
	flterPagination, err := domain.GetPaginationFromCtx(r)
	if err != nil {
		ngelog.Error(r.Context(), "failed to get pagination", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	startDate, _ := utils.ParseDate(r.URL.Query().Get("start_date"))
	endDate, _ := utils.ParseDate(r.URL.Query().Get("end_date"))
	data, pagination, err := h.usecase.GetEvents(r.Context(), domain.EventFilter{
		Title:            r.URL.Query().Get("title"),
		Type:             constants.EventType(r.URL.Query().Get("type")),
		Status:           r.URL.Query().Get("status"),
		StartDate:        startDate,
		EndDate:          endDate,
		FilterPagination: flterPagination,
	})

	if err != nil {
		ngelog.Error(r.Context(), "failed to get events", err)
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:       http.StatusOK,
		Message:    "success",
		Data:       data,
		Pagination: &pagination,
	}, w)
}
