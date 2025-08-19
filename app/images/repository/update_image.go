package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) UpdateImage(ctx context.Context, img domain.Image) error {
	err := repo.db.DB(ctx).Model(&domain.Image{}).
		Where("id = ?", img.ID).
		Updates(map[string]interface{}{
			"file_name":    img.FileName,
			"file_path":    img.FilePath,
			"format":       img.Format,
			"content_type": img.ContentType,
		}).Error
	if err != nil {
		logrus.Error("repo.UpdateImage : failed to update")
		return err
	}
	return nil
}
