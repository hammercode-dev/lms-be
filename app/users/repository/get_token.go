package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) GetToken(ctx context.Context, token string) (logoutToken domain.LogoutToken, err error) {
	if err = repo.db.DB(ctx).Find(&logoutToken, "token = ?", token).Error; err != nil {
		return
	}
	return
}
