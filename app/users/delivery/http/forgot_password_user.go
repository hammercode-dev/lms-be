package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {

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

	_, link, err := h.usecase.ForgotPassword(r.Context(), forgotPassword)

	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Request reset password success check your email",
		Data:    link,
	}, w)

}
