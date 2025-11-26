package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"gorm.io/gorm"
)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) domain.TransactionEventRepository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, data domain.TransactionEvent) (domain.TransactionEvent, error) {
	if err := r.db.WithContext(ctx).Create(&data).Error; err != nil {
		return domain.TransactionEvent{}, err
	}
	return data, nil
}

func (r *repository) GetByTransactionNo(ctx context.Context, transactionNo string) (domain.TransactionEvent, error) {
	var result domain.TransactionEvent
	err := r.db.WithContext(ctx).
		Preload("Registration").
		Preload("Registration.Event").
		Preload("Registration.User").
		Where("transaction_no = ?", transactionNo).
		First(&result).Error
	return result, err
}

func (r *repository) GetByRegistrationID(ctx context.Context, registrationID uint) (domain.TransactionEvent, error) {
	var result domain.TransactionEvent
	err := r.db.WithContext(ctx).
		Where("registration_id = ?", registrationID).
		First(&result).Error
	return result, err
}

func (r *repository) GetByExternalID(ctx context.Context, externalID string) (domain.TransactionEvent, error) {
	var result domain.TransactionEvent
	err := r.db.WithContext(ctx).
		Preload("Registration").
		Where("external_id = ?", externalID).
		First(&result).Error
	return result, err
}

func (r *repository) Update(ctx context.Context, data domain.TransactionEvent) error {
	return r.db.WithContext(ctx).Save(&data).Error
}
