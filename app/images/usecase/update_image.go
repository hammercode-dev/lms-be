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

func (us *usecase) UpdateImage(ctx context.Context, payload domain.UploadImage, fileName string) (domain.UploadImageResponse, error) {
	// Get existing image data by fileName
	existingImage, err := us.imageRepo.GetImage(ctx, fileName)
	if err != nil {
		logrus.Error("failed to get existing image:", err)
		return domain.UploadImageResponse{}, err
	}

	var filePath string
	hName := hash.GenerateHash(time.Now().Format("2006-01-02 15:04:05") + "Hammercode")
	uploadDir := "./uploads"
	if payload.UserID != "" {
		uploadDir = fmt.Sprintf("./uploads/%s", payload.UserID)
	}

	// file path/category/type image
	uploadDir = fmt.Sprintf("%s/%s/%s", uploadDir, payload.Category, payload.Type)
	newFileName := fmt.Sprintf("%s.%s", hName[0:15], payload.Format)
	filePath = fmt.Sprintf("%s/%s", uploadDir, newFileName)

	// Ensure the directory exists, create it if it doesn't
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		logrus.Error("Failed to create directory:", err)
		return domain.UploadImageResponse{}, err
	}

	// Save the file to the uploads directory
	dst, err := os.Create(filePath)
	if err != nil {
		logrus.Error("failed to create file path")
		return domain.UploadImageResponse{}, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, payload.File)
	if err != nil {
		logrus.Error("failed to read all file")
		return domain.UploadImageResponse{}, err
	}

	// Update data di database with existing ID
	img := domain.Image{
		ID:          existingImage.ID,
		FileName:    newFileName,
		FilePath:    filePath,
		Format:      payload.Format,
		ContentType: payload.ContentType,
	}

	err = us.imageRepo.UpdateImage(ctx, img)
	if err != nil {
		return domain.UploadImageResponse{}, err
	}

	return domain.UploadImageResponse{
		FileName: newFileName,
	}, nil
}
