package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	usecase domain.BlogPostUsecase
}

// CreateBlogPost implements domain.BlogPostHandler.
func (h Handler) CreateBlogPost(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	var token string
	if hasToken := strings.HasPrefix(r.Header.Get("Authorization"), "Bearer "); hasToken {
		token = strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if token == "" {
			resp := utils.CustomErrorResponse(utils.NewUnauthorizedError(r.Context(), "Authorization token is required", nil))
			utils.Response(resp, w)
			return
		}
	}

	BlogPost := domain.BlogPost{}
	if err = json.Unmarshal(bodyBytes, &BlogPost); err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	data, err := h.usecase.CreateBlogPost(r.Context(), BlogPost, token)
	if err != nil {
		resp := utils.CustomErrorResponse(err)
		utils.Response(resp, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusCreated,
		Message: "Blog post created successfully",
		Data:    data,
	}, w)

}

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
		PublishedAt time.Time     `json:"published_at"`
		UpdatedAt   time.Time     `json:"updated_at"`
		Tags        []string      `json:"tags" gorm:"-"`
		Category    string        `json:"category"`
		Status      string        `json:"status" gorm:"type:enum('draft', 'published', 'archived')"`
		Slug        string        `json:"slug"`
	}

	responseDTO := []response{}
	for _, post := range data {
		resp := response{
			Id:          post.Id,
			Title:       post.Title,
			Excerpt:     post.Excerpt,
			Author:      post.Author,
			AuthorID:    post.AuthorID,
			PublishedAt: post.PublishedAt,
			UpdatedAt:   post.PublishedAt,
			Tags:        post.Tags,
			Category:    post.Category,
			Status:      post.Status,
			Slug:        post.Slug,
		}
		responseDTO = append(responseDTO, resp)
	}

	utils.Response(domain.HttpResponse{
		Code:       http.StatusOK,
		Message:    "Blog posts retrieved successfully",
		Data:       responseDTO,
		Pagination: paginationResponse,
	}, w)
}

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

// UpdateBlogPost implements domain.BlogPostHandler.
func (h Handler) UpdateBlogPost(w http.ResponseWriter, r *http.Request) {
	idS := mux.Vars(r)["id"]
	id, err := strconv.ParseUint(idS, 10, 32)
	if err != nil {
		logrus.Error("failed to convert string to uint: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid ID format",
		}, w)
		return
	}

	existingPost, err := h.usecase.GetDetailBlogPost(r.Context(), "", uint(id))
	if err != nil {
		logrus.Error("failed to get existing blog post: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusNotFound,
			Message: "Blog post not found",
		}, w)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.Error("failed to read body: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request body",
		}, w)
		return
	}

	var patchData map[string]interface{}
	if err := json.Unmarshal(body, &patchData); err != nil {
		logrus.Error("failed to unmarshal: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
		}, w)
		return
	}

	updatedPost := existingPost
	if title, ok := patchData["title"].(string); ok {
		updatedPost.Title = title
	}
	if content, ok := patchData["content"].(string); ok {
		updatedPost.Content = content
	}
	if excerpt, ok := patchData["excerpt"].(string); ok {
		updatedPost.Excerpt = excerpt
	}
	if category, ok := patchData["category"].(string); ok {
		updatedPost.Category = category
	}
	if status, ok := patchData["status"].(string); ok {
		updatedPost.Status = status
	}
	if tags, ok := patchData["tags"].([]interface{}); ok {
		updatedPost.Tags = make([]string, len(tags))
		for i, tag := range tags {
			updatedPost.Tags[i] = tag.(string)
		}
	}

	if authorData, ok := patchData["author"].(map[string]interface{}); ok {
		if avatar, ok := authorData["avatar"].(string); ok && avatar != "" {
			updatedPost.Author.Avatar = avatar
		}
	}

	updatedPost.UpdatedAt = time.Now()

	err = h.usecase.UpdateBlogPost(r.Context(), updatedPost, uint(id))
	if err != nil {
		logrus.Error("failed to update blog post: ", err)
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: "Failed to update blog post",
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    http.StatusOK,
		Message: "Blog post updated successfully",
		Data:    updatedPost,
	}, w)
}

var (
	handlr *Handler
)

func NewHandler(usecase domain.BlogPostUsecase) domain.BlogPostHandler {
	if handlr == nil {
		handlr = &Handler{
			usecase: usecase,
		}

	}
	return *handlr
}
