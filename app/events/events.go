package events

import (
	handler_http "github.com/hammer-code/lms-be/app/events/delivery/http"
	repository "github.com/hammer-code/lms-be/app/events/repository"
	usecase "github.com/hammer-code/lms-be/app/events/usecase"

	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
)

func InitRepository(db db.DatabaseTransaction) domain.EventRepository {
	return repository.NewRepository(db)
}

func InitUsecase(cfg config.Config, repository domain.EventRepository, imageRepository domain.ImageRepository, dbTX db.DatabaseTransaction) domain.EventUsecase {
	return usecase.NewUsecase(cfg, repository, imageRepository, dbTX)
}

func InitHandler(uc domain.EventUsecase) domain.EventHandler {
	return handler_http.NewHandler(uc)
}
