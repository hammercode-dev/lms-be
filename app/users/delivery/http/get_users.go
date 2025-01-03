package http

import (
	"net/http"
	"strconv"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/jwt"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// GetUsers
// @Summary Get Users
// @Description This endpoint use to get users by filter
// @Tags User
// @Accept json
// @Produce json
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Success 200 {object} []domain.User
// @Router /api/users [get]
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

func (h Handler) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		utils.Response(domain.HttpResponse{
			Code:    401,
			Message: "Not permission",
		}, w)
		return
	}

	claims, err := jwt.ParseToken(authorizationHeader, config.GetConfig().JWT_SECRET_KEY)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	user, err := h.usecase.GetUserById(r.Context(), int8(claims.ID))

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
