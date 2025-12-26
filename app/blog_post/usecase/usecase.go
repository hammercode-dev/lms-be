package usecase

import (
	"github.com/hammer-code/lms-be/domain"
	"github.com/hammer-code/lms-be/pkg/db"
)

type usecase struct {
	repo domain.BlogPostRepository
	dbTX db.DatabaseTransaction
}

var (
	usec *usecase
)

func NewUsecase(repo domain.BlogPostRepository, dbTX db.DatabaseTransaction) domain.BlogPostUsecase {
	if usec == nil {
		usec = &usecase{
			repo: repo,
			dbTX: dbTX,
		}
	}
	return usec
}
