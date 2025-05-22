package usecase

import (
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
	"github.com/hammer-code/lms-be/pkg/jwt"
)

type usecase struct {
	userRepo domain.UserRepository
	dbTX     db.DatabaseTransaction
	jwt      jwt.JWT
	cfg      config.Config
}

var (
	usec *usecase
)

func NewUsecase(cfg config.Config, userRepo domain.UserRepository, dbTX db.DatabaseTransaction, jwt jwt.JWT) domain.UserUsecase {
	if usec == nil {
		usec = &usecase{
			cfg:      cfg,
			userRepo: userRepo,
			dbTX:     dbTX,
			jwt:      jwt,
		}
	}
	return usec
}
