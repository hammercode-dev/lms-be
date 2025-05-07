package users

import (
	users_handler "github.com/hammer-code/lms-be/app/users/delivery/http"
	users_repo "github.com/hammer-code/lms-be/app/users/repository"
	users_usecase "github.com/hammer-code/lms-be/app/users/usecase"
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
	"github.com/hammer-code/lms-be/pkg/jwt"
)

func InitRepository(db db.DatabaseTransaction) domain.UserRepository {
	return users_repo.NewRepository(db)
}

func InitUsecase(cfg config.Config, repository domain.UserRepository, dbTX db.DatabaseTransaction, jwt jwt.JWT) domain.UserUsecase {
	return users_usecase.NewUsecase(cfg, repository, dbTX, jwt)
}

func InitHandler(usecase domain.UserUsecase) domain.UserHandler {
	return users_handler.NewHandler(usecase)
}
