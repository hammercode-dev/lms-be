package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

// set allowed origin
// only request from allowed origin, can bed passed to this endpoint
func (h handler) Subscribe(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	payload := domain.SubscribeReq{}
	if err := json.Unmarshal(body, &payload); err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	err = h.usecase.Subscribe(r.Context(), payload.Email)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Subscribed sucessfully",
	}, w)
}
