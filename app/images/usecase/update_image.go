package usecase

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/hash"
	"github.com/sirupsen/logrus"
)

func (us *usecase) UpdateImage(ctx context.Context, payload domain.UploadImage, id uint) error {
	var filePath string
	hName := hash.GenerateHash(time.Now().Format("2006-01-02 15:04:05") + "Hammercode")
	uploadDir := "./uploads"
	if payload.UserID != "" {
		uploadDir = fmt.Sprintf("./uploads/%s", payload.UserID)
	}

	// file path/category/type image
	uploadDir = fmt.Sprintf("%s/%s/%s", uploadDir, payload.Category, payload.Type)
	fileName := fmt.Sprintf("%s.%s", hName[0:15], payload.Format)
	filePath = fmt.Sprintf("%s/%s", uploadDir, fileName)

	// Ensure the directory exists, create it if it doesn't
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		logrus.Error("Failed to create directory:", err)
		return err
	}

	// Save the file to the uploads directory
	dst, err := os.Create(filePath)
	if err != nil {
		logrus.Error("failed to create file path")
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, payload.File)
	if err != nil {
		logrus.Error("failed to read all file")
		return err
	}

	// Update data di database
	img := domain.Image{
		ID:          id,
		FileName:    fileName,
		FilePath:    filePath,
		Format:      payload.Format,
		ContentType: payload.ContentType,
	}

	return us.imageRepo.UpdateImage(ctx, img)
}
