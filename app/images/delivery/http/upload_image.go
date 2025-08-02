package http

import (
	"fmt"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

func (h Handler) UploadImage(w http.ResponseWriter, r *http.Request) {
	// Limit the size of the uploaded file to 10MB
	r.ParseMultipartForm(10 << 20)

	// Retrieve the file from the form-data
	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving file: %v", err), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Log the file details for debugging
	logrus.Infof("Uploading file: %s, size: %d, MIME header: %v", header.Filename, header.Size, header.Header)

	category := r.FormValue("category")
	if category == "" {
		category = "public"
	}

	contentType := header.Header.Values("Content-Type")[0]
	contentFiles := strings.Split(contentType, "/")

	payload := domain.UploadImage{
		File:        file,
		Header:      header,
		Category:    category,
		ContentType: contentType,
		Format:      contentFiles[1],
		Type:        contentFiles[0],
	}

	resp, err := h.usecase.UploadImage(r.Context(), payload)
	if err != nil {
		utils.Response(domain.HttpResponse{
			Code:    500,
			Message: err.Error(),
		}, w)
		return
	}

	utils.Response(domain.HttpResponse{
		Code:    201,
		Message: "success",
		Data:    resp,
	}, w)
}
