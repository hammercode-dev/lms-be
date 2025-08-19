package domain

import (
	"context"
	"mime/multipart"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type ImageRepository interface {
	Store(context.Context, Image) (uint, error)
	GetImage(ctx context.Context, fileName string) (Image, error)
	UpdateImage(ctx context.Context, file Image) error
	UpdateUseImage(ctx context.Context, id uint) error
}

type ImageUsecase interface {
	UploadImage(context.Context, UploadImage) (UploadImageResponse, error)
	UpdateImage(ctx context.Context, file UploadImage, id uint) error
	UpdateUseImage(context.Context, uint) error
	GetStorage(ctx context.Context, fileName string) (filePath string, err error)
}

type ImageHandler interface {
	UpdateImage(w http.ResponseWriter, r *http.Request)
	UploadImage(w http.ResponseWriter, r *http.Request)
	GetStorage(w http.ResponseWriter, r *http.Request)
}

type UploadImage struct {
	File        multipart.File        `json:"file"`
	Header      *multipart.FileHeader `json:"file_header"`
	Format      string                `json:"format"`
	Type        string                `json:"type"`
	ContentType string                `json:"content_type"`
	Category    string                `json:"category"`
	UserID      string
}

type UploadImageResponse struct {
	FileName string `json:"file_name"`
}

type Image struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	FileName    string    `json:"file_name"`
	FilePath    string    `json:"file_path"`
	Format      string    `json:"format"`
	ContentType string    `json:"content_type"`
	IsUsed      bool      `json:"is_has_been_used"`
	FileSize    int64     `json:"file_size"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   null.Time `json:"updated_at"`
	DeletedAt   null.Time `json:"deleted_at"`
}
