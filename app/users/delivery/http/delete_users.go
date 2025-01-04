package http

import (
	"net/http"
	"strconv"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// DeleteUser
// @Summary Delete User
// @Description This endpoint use to delete user
// @Tags User
// @Accept json
// @Produce json
// @Param id query string false "string"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.HttpResponse
// @Router /api/users [delete]
func (h Handler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)

	err := h.usecase.DeleteUser(r.Context(), int8(userID))

	if err != nil {
		logrus.Error("userUsecase: failed to delete user")
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
	}, w)
}
