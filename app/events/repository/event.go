package repository

import (
	"context"

	"github.com/hammer-code/lms-be/domain"
	"github.com/sirupsen/logrus"
)

func (r repository) DeleteEvent(ctx context.Context, id uint) error {
	err := repo.db.DB(ctx).Model(&domain.Event{}).Delete("id = ?", id).Error
	if err != nil {
		logrus.Error("failed to delete event")
		return err
	}
	return nil
}

func (r repository) CreatePayEvent(ctx context.Context, event domain.EventPay) (uint, error) {
	err := repo.db.DB(ctx).Create(&event).Error
	if err != nil {
		logrus.Error("failed to create event")
		return 0, err
	}
	return event.ID, nil
}

func (r repository) CreateEvent(ctx context.Context, event domain.Event) (uint, error) {
	err := repo.db.DB(ctx).Create(&event).Error
	if err != nil {
		logrus.Error("failed to create event")
		return 0, err
	}
	return event.ID, nil
}

func (r repository) CreateEventTag(ctx context.Context, data domain.EventTag) (uint, error) {
	err := repo.db.DB(ctx).Create(&data).Error
	if err != nil {
		logrus.Error("failed to create event tag")
		return 0, err
	}
	return data.ID, nil
}

func (r repository) CreateEventSpeaker(ctx context.Context, data domain.EventSpeaker) (uint, error) {
	err := repo.db.DB(ctx).Create(&data).Error
	if err != nil {
		logrus.Error("failed to create event tag")
		return 0, err
	}
	return data.ID, nil
}

func (r repository) CreateRegisterEvent(ctx context.Context, event domain.RegistrationEvent) (uint, error) {
	err := repo.db.DB(ctx).Create(&event).Error
	if err != nil {
		logrus.Error("failed to create event")
		return 0, err
	}
	return event.ID, nil
}

// query untuk get user di database
func (repo *repository) GetEvents(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.Event, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	var totalData int64

	if filter.Type != "" {
		db = db.Where("type = ?", filter.Type)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	if filter.StartDate.Valid {
		db = db.Where("start_date > ?", filter.StartDate)
	}

	if filter.StartDate.Valid {
		db = db.Where("end_date < ?", filter.EndDate)
	}

	if filter.Title != "" {
		db = db.Where("title ILIKE ?", "%"+filter.Title+"%")
	}

	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Preload("Tags").Preload("Speakers").Find(&data).Error
	if err != nil {
		logrus.Error("repo.GetEvents: failed to get events use generic conditions")
		return
	}

	return int(totalData), data, err
}

func (repo *repository) GetEvent(ctx context.Context, eventID uint) (data domain.Event, err error) {
	db := repo.db.DB(ctx).Model(&domain.Event{})

	err = db.Where("id = ?", eventID).Find(&data).Error
	if err != nil {
		logrus.Error("repo.GetEvents: failed to get events use generic conditions")
		return
	}

	return data, err
}

func (repo *repository) GetRegistrationEvent(ctx context.Context, orderNo string) (data domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get registration event use generic conditions")
		return
	}

	return data, err
}

func (repo *repository) ListRegistration(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.RegistrationEvent, err error) {
	db := repo.db.DB(ctx).Model(&domain.RegistrationEvent{})

	var totalData int64

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	if filter.StartDate.Valid {
		db = db.Where("start_date > ?", filter.StartDate)
	}

	if filter.StartDate.Valid {
		db = db.Where("end_date < ?", filter.EndDate)
	}

	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).Find(&data).Error
	if err != nil {
		logrus.Error("failed to list registration event use generic conditions")
		return
	}

	return int(totalData), data, err
}

func (repo *repository) GetEventPay(ctx context.Context, orderNo string) (data domain.EventPay, err error) {
	db := repo.db.DB(ctx).Model(&domain.EventPay{})

	err = db.Where("order_no = ?", orderNo).Find(&data).Error
	if err != nil {
		logrus.Error("failed to get event pay")
		return
	}

	return data, err
}

func (repo *repository) ListEventPay(ctx context.Context, filter domain.EventFilter) (tData int, data []domain.EventPay, err error) {
	db := repo.db.DB(ctx).Model(&domain.EventPay{})

	var totalData int64

	if filter.ID != 0 {
		db = db.Where("event_id = ?", filter.ID)
	}

	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}

	if filter.StartDate.Valid {
		db = db.Where("start_date > ?", filter.StartDate)
	}

	if filter.StartDate.Valid {
		db = db.Where("end_date < ?", filter.EndDate)
	}

	db.Count(&totalData)

	err = db.Limit(filter.FilterPagination.GetLimit()).
		Offset(filter.FilterPagination.GetOffset()).
		Preload("RegistrationEvent").Find(&data).Error
	if err != nil {
		logrus.Error("repo.GetEvents: failed to get event pays use generic conditions")
		return
	}

	return int(totalData), data, err
}

func (r repository) UpdateEventPay(ctx context.Context, event domain.EventPay) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update event pay")
		return err
	}
	return nil
}

func (r repository) UpdateRegistrationEvent(ctx context.Context, event domain.RegistrationEvent) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update registration event pay")
		return err
	}
	return nil
}

func (r repository) UpdateEvent(ctx context.Context, event domain.Event) error {
	err := repo.db.DB(ctx).Save(&event).Error
	if err != nil {
		logrus.Error("failed to update event")
		return err
	}

	return nil
}
