package repository

import (
	"context"
	"time"

	"github.com/hammer-code/lms-be/domain"
)

func (repo *repository) StoreToken(ctx context.Context, token string, expiredAt time.Time, uid int) error {
	err := repo.db.DB(ctx).Create(&domain.LogoutToken{
		Token:     token,
		ExpiredAt: expiredAt,
		CreatedAt: time.Now(),
		UserId:    uid,
		Status:    1,
	}).Error

	if err != nil {
		return err
	}

	return nil
}
