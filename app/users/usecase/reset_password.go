package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (us *usecase) ResetPassword(ctx context.Context, reqBodyInstance domain.ForgotPassword) error {
	jwtData, err := us.jwt.VerifyToken(reqBodyInstance.Token)
	if err != nil {
		logrus.Error("us.ResetPassword: failed to verify token", err)
		return err
	}

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(reqBodyInstance.Password), bcrypt.DefaultCost)
	if err != nil {
		logrus.Error("us.ResetPassword: failed to hash password", err)
		return err
	}

	if err := us.userRepo.ResetPassword(ctx, jwtData.Email, string(hashPassword)); err != nil {
		logrus.Error("us.ResetPassword: failed to reset password", err)
		return err
	}

	return nil

}
