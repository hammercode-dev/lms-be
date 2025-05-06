package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) ResetPassword(ctx context.Context, email string, password string) error {
	err := repo.db.DB(ctx).Model(domain.User{}).Where("email = ?", email).Update("password", password).Error

	if err != nil {
		logrus.Error("repo.FindByEmail: failed to find user")
		return err
	}
	return nil
}
