package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)

	}

	forgotPasswordInstance := domain.ForgotPassword{}

	if err := json.Unmarshal(bodyBytes, &forgotPasswordInstance); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
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
