package domain

import (
	"time"

	"gopkg.in/guregu/null.v4"
)

type RegisterEventPayload struct {
	EventID           uint    `json:"event_id"`
	UserID            string  `json:"user_id"`
	Name              string  `json:"name"`
	Email             string  `json:"email"`
	PhoneNumber       string  `json:"phone_number"`
	ImageProofPayment string  `json:"image_proof_payment"`
	NetAmount         float64 `json:"net_amount"`
}

type RegistrationEvent struct {
	ID                uint      `json:"id" gorm:"primarykey"`
	OrderNo           string    `json:"order_no"`
	EventID           uint      `json:"event_id"` // lock event
	UserID            string    `json:"user_id"`  // lock user
	ImageProofPayment string    `json:"image_proof_payment"`
	PaymentDate       null.Time `json:"payment_date"`
	Status            string    `json:"status"` // register, pay, approve/cancel/decline
	UpToYou           string    `json:"-"`
	CreatedByUserID   int       `json:"-"`
	UpdatedByUserID   int       `json:"-"`
	DeletedByUserID   int       `json:"-"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         null.Time `json:"-"`
	DeletedAt         null.Time `json:"-"`
	Event             Event     `json:"event_detail" gorm:"foreignKey:EventID"`
	User              User      `json:"user_detail" gorm:"foreignKey:UserID"`
}

func (RegistrationEvent) TableName() string {
	return "registration_events"
}

type RegistrationEventDTO struct {
	OrderNo string `json:"order_no"`
}

type RegisterEventResponse struct {
	OrderNo string `json:"order_no"`
}

type RegisterStatusResponse struct {
	OrderNo string `json:"order_no"`
	Status  string `json:"string"`
}

type EventRegistrationDTO struct {
	RegistrationID    uint                `json:"registration_id"`
	OrderNo           string              `json:"order_no"`
	PaymentDate       time.Time           `json:"payment_date"`
	Status            string              `json:"status"`
	ImageProofPayment string              `json:"image_proof_payment"`
	UserDetail        RegistrationUserDTO `json:"user_detail"`
}

type RegistrationUserDTO struct {
	UserID      uint   `json:"user_id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Fullname    string `json:"fullname"`
	PhoneNumber string `json:"phone_number"`
}

func (d RegistrationEvent) ToRegistrationDTO() EventRegistrationDTO {
	return EventRegistrationDTO{
		RegistrationID:    d.ID,
		OrderNo:           d.OrderNo,
		PaymentDate:       d.PaymentDate.Time,
		Status:            d.Status,
		ImageProofPayment: d.ImageProofPayment,
		UserDetail: RegistrationUserDTO{
			UserID:      uint(d.User.ID),
			Username:    d.User.Username,
			Email:       d.User.Email,
			Fullname:    d.User.Fullname,
			PhoneNumber: d.User.PhoneNumber,
		},
	}
}

type UpdateRegistrationStatusRequest struct {
	Status string `json:"status"`
}
