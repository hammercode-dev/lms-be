package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

// Login
// @Summary Login
// @Description This endpoint use to login
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body domain.Login true "Body"
// @Failure 400 {object} domain.HttpResponse
// @Failure 500 {object} domain.HttpResponse
// @Failure 200 {object} domain.HttpResponse
// @Router /api/v1/auth/login [post]
func (h Handler) Login(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	loginInstance := domain.Login{}
	if err := json.Unmarshal(bodyBytes, &loginInstance); err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	_, token, err := h.usecase.Login(r.Context(), loginInstance)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Login successfully",
		Data:    token,
	}, w)
}
