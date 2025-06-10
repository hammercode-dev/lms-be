package usecase

import (
	"context"

	"github.com/hammer-code/lms-be/constants"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/utils"
	"golang.org/x/crypto/bcrypt"
)

func (us *usecase) Register(ctx context.Context, userReq domain.User) (domain.User, error) {
	user := domain.User{}
	err := us.dbTX.StartTransaction(ctx, func(txCtx context.Context) error {
		hashPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}

		userReq.Password = string(hashPassword)
		userReq.Role = constants.RoleUser
		user, err = us.userRepo.CreateUser(ctx, userReq)
		if err != nil {
			err = utils.NewInternalServerError(ctx, err)
			return err
		}
		return nil

	})

	if err != nil {
		err = utils.NewInternalServerError(ctx, err)
		return domain.User{}, err
	}
	
	return user, nil
}
