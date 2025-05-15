package repository

import (
	"context"
	"time"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) ForgotPassword(ctx context.Context, token string, expiredAt time.Time, user domain.User) (err error) {
	err = repo.db.DB(ctx).Create(&domain.ResetPasswordToken{
		Token:      token,
		UserID:     uint64(user.ID),
		ExpiryDate: expiredAt,
	}).Error
	if err != nil {
		logrus.Error("repo.ForgotPassword : failed to create reset password token")
		return err
	}

	return nil
}
