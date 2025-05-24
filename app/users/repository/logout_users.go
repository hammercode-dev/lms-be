package repository

import (
	"context"
	"time"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) LogoutUser(ctx context.Context, token string, expiredAt time.Time) error {
	if err := repo.db.DB(ctx).Model(&domain.LogoutToken{}).Where("token = ?", token).Update("status", 0).Error; err != nil {
		return err
	}

	return nil
}
