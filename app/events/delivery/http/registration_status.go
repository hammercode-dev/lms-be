package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
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
		logrus.Error("failed to Create pay event : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}, w)
}
