package http

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	reEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		}, w)
		return
	}

	forgotPassword := domain.ForgotPassword{}
	if err = json.Unmarshal(bodyBytes, &forgotPassword); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	if isValidEmail := reEmail.MatchString(forgotPassword.Email); !isValidEmail {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Email is not valid",
		}, w)
		return
	}

	if err := h.usecase.ForgotPassword(r.Context(), forgotPassword); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Request reset password success check your email",
		Data:    nil,
	}, w)

}
