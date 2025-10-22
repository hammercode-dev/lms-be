package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

// NewRepository creates simple repository
func NewRepository(db *gorm.DB) domain.TestingTransactionRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data domain.TestingTransaction) error {
	return r.db.WithContext(ctx).Create(&data).Error
}

func (r *repository) GetByOrderNo(ctx context.Context, orderNo string) (domain.TestingTransaction, error) {
	var data domain.TestingTransaction
	err := r.db.WithContext(ctx).Where("order_no = ?", orderNo).First(&data).Error
	return data, err
}

func (r *repository) GetAll(ctx context.Context) ([]domain.TestingTransaction, error) {
	var data []domain.TestingTransaction
	err := r.db.WithContext(ctx).Order("created_at DESC").Find(&data).Error
	return data, err
}

func (r *repository) Update(ctx context.Context, data domain.TestingTransaction) error {
	return r.db.WithContext(ctx).Save(&data).Error
}
