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
	
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(userReq.Password)); err != nil {
			logrus.Error("us.Login: invalid password")
			return err
		}
	
		
		tokenPtr, expiredTime, err := us.jwt.GenerateAccessToken(ctx, &user, 60)
		token = *tokenPtr
		
		if err != nil {
			logrus.Error("us.Login: failed to login. ", err)
			return err
		}

		if err = us.userRepo.UnactivateTokenByUser(ctx, user.ID); err != nil {
			logrus.Error("us.Login: failed to login. ", err)
			return err
		}
		if err = us.userRepo.StoreToken(ctx, token, expiredTime, user.ID); err != nil {
			logrus.Error("us.Login: failed to login. ", err)
			return err
		}

		return nil
	})

	if err != nil {
		logrus.Error("us.Login: failed to login. ", err)
		return
	}

	return user, token, nil
}
