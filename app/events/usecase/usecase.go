package usecase

import (
	"github.com/hammer-code/lms-be/config"
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
	"github.com/hammer-code/lms-be/pkg/jwt"
)

type usecase struct {
	repository      domain.EventRepository
	imageRepository domain.ImageRepository
	cfg             config.Config
	dbTX            db.DatabaseTransaction
	jwt             jwt.JWT
}

var (
	uc *usecase
)

func NewUsecase(cfg config.Config, repository domain.EventRepository, imageRepository domain.ImageRepository, dbTX db.DatabaseTransaction, jwt jwt.JWT) domain.EventUsecase {
	if uc == nil {
		uc = &usecase{
			repository:      repository,
			imageRepository: imageRepository,
			dbTX:            dbTX,
			cfg:             cfg,
			jwt:             jwt,
		}
	}

	return uc
}
