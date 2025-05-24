package repository

import (
	"context"
	"time"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) CleanupLogoutToken(ctx context.Context) error {
	token := &domain.LogoutToken{}
	if err := repo.db.DB(ctx).Delete(token, "expired_at < ?", time.Now()).Error; err != nil {
		return err
	}

	return nil
}
