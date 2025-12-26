package http

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

// GetDetailBlogPost implements domain.BlogPostHandler.
func (h Handler) GetDetailBlogPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	slug := vars["slug"]

	if slug == "" {
		resp := utils.CustomErrorResponse(utils.NewBadRequestError(r.Context(), "Slug is required", nil))
		utils.Response(resp, w)
		return
	}

	resp, err := h.usecase.GetDetailBlogPost(r.Context(), slug, 0)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "Blog post retrieved successfully",
		Data:    resp,
	}, w)
}
