package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) UnactivateTokenByUser(ctx context.Context, uid int) error {
	if err := repo.db.DB(ctx).Model(&domain.LogoutToken{}).Where("user_id = ? AND status = 1", uid).Update("status", 0).Error; err != nil {
		return err
	}

	return nil
}
