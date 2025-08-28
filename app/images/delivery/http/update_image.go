package http

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	// Parse image ID from query or URL (misal: /images/{id})
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "missing image id",
		}, w)
		return
	}
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "invalid image id",
		}, w)
		return
	}

	// Parse multipart form
	err = r.ParseMultipartForm(10 << 20) // 10MB
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to parse form",
		}, w)
		return
	}
	file, header, err := r.FormFile("image")
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "failed to get file",
		}, w)
		return
	}
	defer file.Close()

	category := r.FormValue("category")
	if category == "" {
		category = "public"
	}

	contentType := header.Header.Values("Content-Type")[0]
	contentFiles := strings.Split(contentType, "/")

	upload := domain.UploadImage{
		File:        file,
		Header:      header,
		Category:    category,
		ContentType: contentType,
		Format:      contentFiles[1],
		Type:        contentFiles[0],
	}

	ctx := r.Context()
	err = h.usecase.UpdateImage(ctx, upload, uint(id))
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusInternalServerError,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    200,
		Message: "Image updated successfully",
	}, w)
}

