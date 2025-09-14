package http

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
)

func (h Handler) UpdateImage(w http.ResponseWriter, r *http.Request) {
	// Parse image fileName from URL (misal: /images/{fileName})
	fileName := mux.Vars(r)["fileName"]
	if fileName == "" {
		utils.Response(domain.HttpResponse{
			Code:    http.StatusBadRequest,
			Message: "missing image file name",
		}, w)
		return
	}

	// Parse multipart form
	err := r.ParseMultipartForm(10 << 20) // 10MB
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
	response, err := h.usecase.UpdateImage(ctx, upload, fileName)
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
		Data:    response,
	}, w)
}

