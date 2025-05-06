package http

import (
	"net/http"

	"html/template"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) ShowResetPasswordForm(w http.ResponseWriter, r *http.Request) {
	// get token from query
	token := r.URL.Query().Get("token")
	if token == "" {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Token is required",
		}, w)
		return
	}

	if err := h.usecase.VerifyPasswordResetToken(r.Context(), token); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    400,
			Message: "Invalid token",
		}, w)
		return
	}

	html, err := template.ParseFiles("./assets/reset-password.html")
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: "Failed to parse template",
			Data:    nil,
		}, w)
		return
	}

	if err := html.Execute(w, token); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: "Failed to execute template",
		}, w)
		return
	}
}
