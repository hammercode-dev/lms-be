package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (us *usecase) ForgotPassword(ctx context.Context, emailForgot domain.ForgotPassword) (user domain.User, resetLink string, err error) {
	err = us.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		user, err = us.userRepo.FindByEmail(ctx, emailForgot.Email)
		if err != nil {
			logrus.Error("us.ForgotPassword: failed to get Email", err)
			return err
		}
		return nil
	})

	if err != nil {
		logrus.Error("us.ForgotPassword: failed to get Email", err)
		return
	}

	resetToken, err := us.jwt.GenerateAccessToken(ctx, &user)
	if err != nil {
		logrus.Error("us.ForgotPassword: failed to generate token", err)
		return
	}
	// link to reset password
	link := "http://localhost:8000/api/v1/auth/forgot_password?token=" + *resetToken

	return user, link, nil
}
