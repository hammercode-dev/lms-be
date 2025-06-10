package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/utils"
)

func (us *usecase) Logout(ctx context.Context, token string) error {
	jwtData, err := us.jwt.VerifyToken(token)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}

	err = us.userRepo.LogoutUser(ctx, token, jwtData.ExpiresAt.Time)
	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return err
	}
	return err
}
