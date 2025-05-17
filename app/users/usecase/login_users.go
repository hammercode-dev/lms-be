package usecase

import (
	"context"
	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (us *usecase) Login(ctx context.Context, userReq domain.Login) (user domain.User, token string, err error) {
	err = us.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		user, err = us.userRepo.FindByEmail(ctx, userReq.Email)
		if err != nil {
			logrus.Error("us.LoginUser: failed to login", err)
			return err
		}
		return nil
	})

	if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password)); err != nil {
		logrus.Error("us.Login: invalid password")
		return
	}

	if err != nil {
		logrus.Error("us.Login: failed to login. ", err)
		return
	}

	signToken, err := us.jwt.GenerateAccessToken(ctx, &user, 60)
	if err != nil {
		logrus.Error("us.Login: failed to login. ", err)
		return
	}

	return user, *signToken, nil
}
