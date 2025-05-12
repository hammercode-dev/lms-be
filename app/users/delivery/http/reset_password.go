package http

import (
	"encoding/json"
	"io"
	"net/http"
	"regexp"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {

	reEmail := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`)
	passwordRegex := regexp.MustCompile(`^[a-zA-Z\d]{8,}$`)

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	forgotPasswordInstance := domain.ForgotPassword{}

	if err := json.Unmarshal(bodyBytes, &forgotPasswordInstance); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	if isValidEmail := reEmail.MatchString(forgotPasswordInstance.Email); isValidEmail != true {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Email is not valid",
		}, w)
		return
	}

	if forgotPasswordInstance.Password != forgotPasswordInstance.ConfirmPassword {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Password and confirm password must be the same",
		}, w)
		return
	}

	if isValidPass := passwordRegex.MatchString(forgotPasswordInstance.Password); isValidPass != true {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Password must contain at least 8 characters, one uppercase letter, one lowercase letter, and one number",
		}, w)
		return
	}

	if err := h.usecase.ResetPassword(r.Context(), forgotPasswordInstance); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return

	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Reset password success",
	}, w)
}
