package http

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// DeleteBlogPost implements domain.BlogPostHandler.
func (h Handler) DeleteBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	idString := vars["id"]

	value, err := strconv.ParseUint(idString, 10, 32)
	if err != nil {
		logrus.Error("failed to convert string to uint: ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	err = h.usecase.DeleteBlogPost(r.Context(), uint(value))
	if err != nil {
		logrus.Error("failed to delete event : ", err)
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "success",
		Data:    nil,
	}, w)
}
