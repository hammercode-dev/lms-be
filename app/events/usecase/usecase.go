package usecase

import (
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
)

type usecase struct {
	repository      domain.EventRepository
	imageRepository domain.ImageRepository
	cfg             config.Config
	dbTX            db.DatabaseTransaction
}

var (
	uc *usecase
)

func NewUsecase(cfg config.Config, repository domain.EventRepository, imageRepository domain.ImageRepository, dbTX db.DatabaseTransaction) domain.EventUsecase {
	if uc == nil {
		uc = &usecase{
			repository:      repository,
			imageRepository: imageRepository,
			dbTX:            dbTX,
			cfg:             cfg,
		}
	}

	return uc
}
