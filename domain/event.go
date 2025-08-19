package domain

import (
	"context"
	"net/http"
	"time"

	"gopkg.in/guregu/null.v4"
)

type EventRepository interface {
	CreateEvent(ctx context.Context, data Event) (uint, error)
	// UpdateEvent(ctx context.Context, payload Event) error
	// DeleteEvent(ctx context.Context, id int) error
	// GetEventByID(ctx context.Context, id int) (Event, error)
	// GetEvents(ctx context.Context, filter EventFilter) ([]Event, error)
	CreateEventTag(ctx context.Context, data EventTag) (uint, error)
	CreateEventSpeaker(ctx context.Context, data EventSpeaker) (uint, error)
	GetEvents(ctx context.Context, filter EventFilter) (tData int, data []EventDTO, err error)
	CreateEventPay(ctx context.Context, event EventPay) (uint, error)
	CreateRegistrationEvent(ctx context.Context, event RegistrationEvent) (uint, error)
	GetEvent(ctx context.Context, eventID uint) (data EventDTO, err error)
	DeleteEvent(ctx context.Context, eventID uint) (err error)
	GetRegistrationEvent(ctx context.Context, orderNo string) (data RegistrationEvent, err error)
	ListRegistration(ctx context.Context, filter EventFilter, email string) (tData int, data []RegistrationEvent, err error)
	ListEventPay(ctx context.Context, filter EventFilter) (tData int, data []EventPay, err error)
	UpdateEventPay(ctx context.Context, event EventPay) error
	GetEventPay(ctx context.Context, orderNo string) (data EventPay, err error)
	UpdateRegistrationEvent(ctx context.Context, event RegistrationEvent) error
	UpdateEvent(ctx context.Context, event Event) error
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, payload CreateEventPayload) error
	UpdateEvent(ctx context.Context, id uint, payload UpdateEventPayload) error
	GetEvents(ctx context.Context, filter EventFilter) (data []EventDTO, pagination Pagination, err error)
	CreateRegistrationEvent(ctx context.Context, payload RegisterEventPayload) (RegisterEventResponse, error)
	CreateEventPay(ctx context.Context, payload EventPayPayload) error
	GetEventByID(ctx context.Context, id uint) (resp EventDTO, err error)
	DeleteEvent(ctx context.Context, id uint) (err error)
	RegistrationStatus(ctx context.Context, orderNo string) (resp RegisterStatusResponse, err error)
	ListRegistration(ctx context.Context, filter EventFilter) (resp []RegistrationEvent, pagination Pagination, err error)
	ListEventPay(ctx context.Context, filter EventFilter) (data []EventPay, pagination Pagination, err error)
	PayProcess(ctx context.Context, payload PayProcessPayload) error
}

type EventHandler interface {
	CreateEvent(w http.ResponseWriter, r *http.Request)
	GetEvents(w http.ResponseWriter, r *http.Request)
	List(w http.ResponseWriter, r *http.Request)
	DeleteEvent(w http.ResponseWriter, r *http.Request)
	RegisterEvent(w http.ResponseWriter, r *http.Request)
	PayEvent(w http.ResponseWriter, r *http.Request)
	GetEventByID(w http.ResponseWriter, r *http.Request)
	GetDetail(w http.ResponseWriter, r *http.Request)
	RegistrationStatus(w http.ResponseWriter, r *http.Request)
	ListRegistration(w http.ResponseWriter, r *http.Request)
	ListEventPay(w http.ResponseWriter, r *http.Request)
	PayProcess(w http.ResponseWriter, r *http.Request)
}

type Event struct {
	ID                   uint           `json:"id" gorm:"primarykey"`
	Title                string         `json:"title" `
	Description          string         `json:"description"`
	Slug                 string         `json:"slug"`
	AuthorID             int            `json:"author_id"`
	Author               User           `gorm:"foreignKey:AuthorID;references:ID"` // Ensure foreign key is correctly referenced
	Image                string         `json:"image"`
	Date                 null.Time      `json:"date"`
	Type                 string         `json:"type"`
	Location             string         `json:"location"`
	Duration             string         `json:"duration"`
	Capacity             int            `json:"capacity"`
	Status               string         `json:"status"`                                          // Conference, Tech Talk, Workshop, Webinar, etc.
	Tags                 []EventTag     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"` // Ensure foreign key is correctly referenced
	Speakers             []EventSpeaker `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"` // Ensure foreign key is correctly referenced
	SessionType          string         `json:"session_type"`                                    // online, offline, hybrid
	RegistrationLink     string         `json:"registration_link"`
	Price                float64        `json:"price"` // 0 == free
	ReservationStartDate null.Time      `json:"reservation_start_date"`
	ReservationEndDate   null.Time      `json:"reservation_end_date"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            null.Time      `json:"updated_at"`
	DeletedAt            null.Time      `json:"deleted_at"`
	AdditionalLink       string         `json:"additional_link"`
}

func (Event) TableName() string {
	return "events"
}

type EventTag struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	EventID uint   `json:"event_id"`
	Tag     string `json:"tags"`
}

func (EventTag) TableName() string {
	return "event_tags"
}

type EventSpeaker struct {
	ID      uint   `json:"id" gorm:"primarykey"`
	EventID uint   `json:"event_id"`
	Name    string `json:"name"`
}

func (EventSpeaker) TableName() string {
	return "event_speakers"
}

type CreateEventPayload struct {
	Title                string    `json:"title" validate:"required"`
	Description          string    `json:"description" validate:"required"`
	Author               string    `json:"author" validate:"required"`
	FileName             string    `json:"file_name" validate:"required"`
	Slug                 string    `json:"slug" validate:"required"`
	IsOnline             string    `json:"is_online" validate:"required"`
	Date                 null.Time `json:"date" validate:"required"`
	Type                 string    `json:"type" validate:"required"`
	Location             string    `json:"location" validate:"required"`
	Duration             string    `json:"duration" validate:"required"`
	Status               string    `json:"status" validate:"required"`
	Capacity             int       `json:"capacity" validate:"required"`
	Price                float64   `json:"price"`
	RegistrationLink     string    `json:"registration_link"`
	Tags                 []string  `json:"tags"`
	Speakers             []string  `json:"speakers"`
	ReservationStartDate null.Time `json:"reservation_start_date"`
	ReservationEndDate   null.Time `json:"reservation_end_date"`
	AdditionalLink       string    `json:"additional_link"`
}

type UpdateEventPayload struct {
	Title                string    `json:"title" validate:"required"`
	Description          string    `json:"description" validate:"required"`
	Author               string    `json:"author" validate:"required"`
	FileName             string    `json:"file_name" validate:"required"`
	Slug                 string    `json:"slug" validate:"required"`
	IsOnline             string    `json:"is_online" validate:"required"`
	Date                 null.Time `json:"date" validate:"required"`
	Type                 string    `json:"type" validate:"required"`
	Location             string    `json:"location" validate:"required"`
	Duration             string    `json:"duration" validate:"required"`
	Status               string    `json:"status" validate:"required"`
	Capacity             int       `json:"capacity" validate:"required"`
	Price                float64   `json:"price"`
	RegistrationLink     string    `json:"registration_link"`
	Tags                 []string  `json:"tags"`
	Speakers             []string  `json:"speakers"`
	ReservationStartDate null.Time `json:"reservation_start_date"`
	ReservationEndDate   null.Time `json:"reseveration_end_date"`
	AdditionalLink       string    `json:"additional_link"`
}

type EventDTO struct {
	ID                   uint           `json:"id" gorm:"primarykey"`
	Title                string         `json:"title" `
	Description          string         `json:"description"`
	Slug                 string         `json:"slug"`
	Author               string         `json:"author"`
	Image                string         `json:"image"`
	Date                 null.Time      `json:"date"`
	Type                 string         `json:"type"`
	Location             string         `json:"location"`
	Duration             string         `json:"duration"`
	Capacity             int            `json:"capacity"`
	Status               string         `json:"status"`                                          // Conference, Tech Talk, Workshop, Webinar, etc.
	Tags                 []EventTag     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"` // Ensure foreign key is correctly referenced
	Speakers             []EventSpeaker `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"` // Ensure foreign key is correctly referenced
	SessionType          string         `json:"session_type"`                                    // online, offline, hybrid
	RegistrationLink     string         `json:"registration_link"`
	Price                float64        `json:"price"` // 0 == free
	ReservationStartDate null.Time      `json:"reservation_start_date"`
	ReservationEndDate   null.Time      `json:"reservation_end_date"`
	CreatedAt            time.Time      `json:"created_at"`
	UpdatedAt            null.Time      `json:"updated_at"`
	DeletedAt            null.Time      `json:"deleted_at"`
	AdditionalLink       string         `json:"additional_link"`
}

type UpdateEvenPayload struct {
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Author           string    `json:"author"`
	ImageEvent       string    `json:"image_event"`
	DateEvent        null.Time `json:"date_event"`
	Type             string    `json:"type"`
	Location         string    `json:"location"`
	Duration         string    `json:"duration"`
	Capacity         int       `json:"capacity"`
	RegistrationLink string    `json:"registration_link"`
}

type EventFilter struct {
	ID        uint
	Title     string
	Type      string
	Status    string
	StartDate null.Time
	EndDate   null.Time
	FilterPagination
}

type EventPay struct {
	ID                  uint              `json:"id" gorm:"primarykey"`
	RegistrationEventID uint              `json:"registration_event_id"`
	EventID             uint              `json:"event_id"`
	ImageProofPayment   string            `json:"image_proof_payment"`
	OrderNO             string            `json:"order_no"`
	NetAmount           float64           `json:"net_amount"`
	Status              string            `json:"status"`
	RegistrationEvent   RegistrationEvent `json:"registration_event" gorm:"foreignKey:RegistrationEventID"`
}

func (EventPay) TableName() string {
	return "event_pays"
}

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
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         null.Time `json:"updated_at"`
	DeletedAt         null.Time `json:"deleted_at"`
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

type EventPayPayload struct {
	OrderNo           string  `json:"order_no"`
	ImageProofPayment string  `json:"image_proof_payment"`
	NetAmount         float64 `json:"net_amount"`
}

type RegisterStatusResponse struct {
	OrderNo string `json:"order_no"`
	Status  string `json:"string"`
}

type PayProcessPayload struct {
	OrderNo string `json:"order_no"`
	Status  string `json:"status"`
	Note    string `json:"note"`
}
