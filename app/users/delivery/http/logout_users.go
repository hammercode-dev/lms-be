package http

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"net/http"
)

// Logout
// @Summary Logout
// @Description This endpoint use to logout
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.Login true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Failure 200 {object} domain.HttpResponse
// @Router /api/auth/logout [post]
func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	token := utils.ExtractBearerToken(r)

	err := h.usecase.Logout(r.Context(), *token)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
			Data:    nil,
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "successfuly logged out",
		Data: map[string]string{
			"token": *token,
		},
	}, w)
}
