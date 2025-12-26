package http

import (
	"net/http"
	"strconv"

	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
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
// @Router /api/v1/users [delete]
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

// GetUsers
// @Summary Get Users
// @Description This endpoint use to get users by filter
// @Tags User
// @Accept json
// @Produce json
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} []domain.User
// @Router /api/v1/users [get]
func (h Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.usecase.GetUsers(r.Context())

	if err != nil {
		logrus.Error("userUsecase: failed to get users")
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
		Data:    users,
	}, w)
}

func (h Handler) GetUserById(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("id")
	userID, _ := strconv.ParseUint(userIDStr, 10, 64)
	user, err := h.usecase.GetUserById(r.Context(), int8(userID))

	if err != nil {
		logrus.Error("userUsecase: failed to get user")
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
		Data:    user,
	}, w)
}

// @Summary Get User Profile
// @Description This endpoint use to get current user
// @Tags User
// @Accept json
// @Produce json
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} domain.User
// @Router /api/v1/users [get]
func (h Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	userData := r.Context().Value(contextkey.UserKey).(domain.User)

	user, err := h.usecase.GetUserById(r.Context(), int8(userData.ID))

	if err != nil {
		logrus.Error("userUsecase: failed to get user")
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
		Data:    user,
	}, w)
}
