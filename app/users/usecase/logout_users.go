package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/utils"
	"github.com/sirupsen/logrus"
)

func (us *usecase) Logout(ctx context.Context, token string) error {
	jwtData, err := us.jwt.VerifyToken(token)
	if err != nil {
		logrus.Error("us.LogoutUser: failed to varify token", err)
		err = utils.NewInternalServerError(err)
		return err
	}

	err = us.userRepo.LogoutUser(ctx, token, jwtData.ExpiresAt.Time)
	if err != nil {
		logrus.Error("us.LogoutUser: failed to logout", err)
		err = utils.NewInternalServerError(err)
		return err
	}
	return err
}
