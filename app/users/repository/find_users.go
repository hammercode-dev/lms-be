package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) FindByEmail(ctx context.Context, email string) (user domain.User, err error) {
	err = repo.db.DB(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		logrus.Error("repo.FindByEmail: failed to find user")
		return
	}
	return user, nil
}

func (repo *repository) FindById(ctx context.Context, id int8) (user domain.User, err error) {
	err = repo.db.DB(ctx).Select("id", "username", "email", "role", "fullname", "date_of_birth", "gender", "phone_number", "address", "github", "linkedin", "personal_web").Where("id = ?", id).Take(&user).Error
	if err != nil {
		logrus.Error("repo.FindById: failed to find user")
		return
	}
	return user, nil
}
