package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (repo *repository) ListRegistrationByEvent(ctx context.Context, eventID uint, filterPagination domain.FilterPagination) (data []domain.RegistrationEvent, totalCount int64, err error) {
	if err := repo.db.DB(ctx).Model(&domain.RegistrationEvent{}).Where("event_id = ?", eventID).Count(&totalCount).Error; err != nil {
		logrus.Error("failed to count blog posts: ", err)
		return nil, 0, err
	}

	offset := filterPagination.GetOffset()
	limit := filterPagination.GetLimit()
	orderBy := filterPagination.GetOrderBy()

	db := repo.db.DB(ctx).
		Model(&domain.RegistrationEvent{}).
		Where("event_id = ?", eventID).
		Preload("User")

	if orderBy == "" {
		orderBy = "created_at DESC"
	}

	db = db.Order(orderBy)

	if limit > 0 {
		db = db.Limit(limit)
	}
	if offset > 0 {
		db = db.Offset(offset)
	}

	if err := db.Find(&data).Error; err != nil {
		return nil, 0, err
	}

	return data, totalCount, nil
}
