package repository

import (
	"context"
	"errors"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) ResetPassword(ctx context.Context, email, password, token string) error {
	if err := repo.db.StartTransaction(ctx, func(ctx context.Context) error {
		ResetPasswordTokenInstance := domain.ResetPasswordToken{}
		if err := repo.db.DB(ctx).Model(ResetPasswordTokenInstance).Where("token = ?", token).First(&ResetPasswordTokenInstance).Error; err != nil {
			logrus.Error("repo.ResetPassword: failed to find token")
			return err
		}

		if ResetPasswordTokenInstance.IsUsed {
			logrus.Error("repo.ResetPassword: token already used")
			return errors.New("token already used")
		}

		if err := repo.db.DB(ctx).Model(domain.User{}).Where("email = ?", email).Update("password", password).Error; err != nil {
			logrus.Error("repo.ResetPassword: failed to update password")
			return err
		}

		if err := repo.db.DB(ctx).Model(domain.ResetPasswordToken{}).Where("token = ?", token).Update("is_used", true).Error; err != nil {
			logrus.Error("repo.ResetPassword: failed to update token state")
			return err
		}
		return nil
	}); err != nil {
		logrus.Error("repo.ResetPassword: failed to reset password")
		return err
	}
	return nil
}
