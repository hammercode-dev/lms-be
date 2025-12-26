package repository

import (
	"github.com/hammer-code/lms-be/domain"
	pkgDB "github.com/hammer-code/lms-be/pkg/db"
)

type repository struct {
	db pkgDB.DatabaseTransaction
}

func NewRepository(db pkgDB.DatabaseTransaction) domain.BlogPostRepository {
	return &repository{
		db: db,
	}
}
