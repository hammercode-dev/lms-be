package http

import (
	"net/http"
	"time"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/ngelog"
	"github.com/hammer-code/lms-be/utils"
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
// @Router /api/v1/auth/logout [post]
func (h Handler) Logout(w http.ResponseWriter, r *http.Request) {
	// token := utils.ExtractBearerToken(r)

	tokenRaw, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			ngelog.Error(r.Context(), "failed to get token from cookie", nil)
			utils.Response(domain.HttpResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    nil,
			}, w)
			return
		} else {
			utils.Response(domain.HttpResponse{
				Code:    401,
				Message: "Unauthorized",
				Data:    nil,
			}, w)
			return
		}
	}
	token := &tokenRaw.Value
	
	err = h.usecase.Logout(r.Context(), *token)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}
	
	expiredTime := time.Now().Local().Add(time.Duration(0) * time.Minute)
	cookie := http.Cookie{
		Name:    "token",
		Value:   "",
		Expires: expiredTime,
		Path:    "/",
		HttpOnly: true,
		Secure:   config.GetConfig().APP_ENV != "development",
	}

	http.SetCookie(w, &cookie)

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "successfuly logged out",
		Data: map[string]string{
			"token": *token,
		},
	}, w)
}
