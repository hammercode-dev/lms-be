package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
)

// RegistrationStatus
// @Summary Register Status
// @Description This endpoint use to check registration status
// @Tags Event
// @Accept json
// @Produce json
// @Param order_no path string true "ABCXX"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.RegisterStatusResponse
// @Router /api/v1/events/registartions/:order_no [get]
func (h Handler) RegistrationStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	order_no := vars["order_no"]
	data, err := h.usecase.RegistrationStatus(r.Context(), order_no)
	if err != nil {
		ngelog.Error(r.Context(), "failed to get registration status", err)
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
