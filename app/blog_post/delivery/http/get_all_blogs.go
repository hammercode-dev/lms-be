package http

import (
	"net/http"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

// GetAllBlogPosts implements domain.BlogPostHandler.
func (h Handler) GetAllBlogPosts(w http.ResponseWriter, r *http.Request) {
	// Ambil parameter pagination dari request
	pagination, err := domain.GetPaginationFromCtx(r)
	if err != nil {
		logrus.Error("failed to parse pagination parameters: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid pagination parameters",
		}, w)
		return
	}

	// Panggil usecase dengan parameter pagination
	data, paginationResponse, err := h.usecase.GetAllBlogPosts(r.Context(), pagination)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	type response struct {
		Id          int           `json:"id" gorm:"primaryKey"`
		Title       string        `json:"title"`
		Excerpt     string        `json:"excerpt"`
		Author      domain.Author `json:"author" gorm:"foreignKey:AuthorID;references:UserId"`
		AuthorID    int           `json:"author_id" gorm:"column:author_id"`
		Tags        []string      `json:"tags" gorm:"-"`
		Category    string        `json:"category"`
		Status      string        `json:"status" gorm:"type:enum('draft', 'published', 'archived')"`
		Slug        string        `json:"slug"`
		PublishedAt *time.Time    `json:"published_at"`
		UpdatedAt   *time.Time    `json:"updated_at"`
		CreatedAt   *time.Time    `json:"created_at"`
	}

	responseDTO := []response{}
	for _, post := range data {
		resp := response{
			Id:          post.Id,
			Title:       post.Title,
			Excerpt:     post.Excerpt,
			Author:      post.Author,
			AuthorID:    post.AuthorID,
			Tags:        post.Tags,
			Category:    post.Category,
			Status:      post.Status,
			Slug:        post.Slug,
			PublishedAt: post.PublishedAt,
			UpdatedAt:   post.UpdatedAt,
			CreatedAt:   post.CreatedAt,
		}
		responseDTO = append(responseDTO, resp)
	}

	utils.Response(domain.HttpResponse{
		Code:       http.StatusOK,
		Message:    "Blog posts retrieved successfully",
		Data:       responseDTO,
		Pagination: &paginationResponse,
	}, w)
}
