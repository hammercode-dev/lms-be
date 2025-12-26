package http

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hammer-code/lms-be/domain"
	contextkey "github.com/hammer-code/lms-be/pkg/context_key"
	"github.com/hammer-code/lms-be/utils"
)

// CreateBlogPost implements domain.BlogPostHandler.
func (h Handler) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	user := r.Context().Value(contextkey.UserKey).(domain.User)

	BlogPost := domain.BlogPost{}
	if err = json.Unmarshal(bodyBytes, &BlogPost); err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	err = h.usecase.CreateBlogPost(r.Context(), BlogPost, user)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "Blog post created successfully",
	}, w)

}
