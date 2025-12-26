package http

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

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
	if patchData["status"] == "published" {
		if updatedPost.PublishedAt == nil {
			timeNow := time.Now()
			updatedPost.PublishedAt = &timeNow
		}
	} else {
		updatedPost.PublishedAt = nil
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

	timeNow := time.Now()
	updatedPost.UpdatedAt = &timeNow

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
	}, w)
}
