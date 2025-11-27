package domain

import (
	"context"
	"net/http"
	"time"

	"github.com/hammer-code/lms-be/constants"
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
	GetEvents(ctx context.Context, filter EventFilter) (tData int, data []Event, err error)
	CreateEventPay(ctx context.Context, event EventPay) (uint, error)
	CreateRegistrationEvent(ctx context.Context, event RegistrationEvent) (uint, error)
	GetEvent(ctx context.Context, eventID uint) (data Event, err error)
	DeleteEvent(ctx context.Context, eventID uint) (err error)
	GetRegistrationEvent(ctx context.Context, orderNo string) (data RegistrationEvent, err error)
	ListRegistration(ctx context.Context, filter EventFilter, email string) (tData int, data []RegistrationEvent, err error)
	ListEventPay(ctx context.Context, filter EventFilter) (tData int, data []EventPay, err error)
	UpdateEventPay(ctx context.Context, event EventPay) error
	GetEventPay(ctx context.Context, orderNo string) (data EventPay, err error)
	UpdateRegistrationEvent(ctx context.Context, event RegistrationEvent) error
	UpdateEvent(ctx context.Context, event Event) error
	ListRegistrationByEvent(ctx context.Context, eventID uint, pagination FilterPagination) (data []RegistrationEvent, totalCount int64, err error)
	GetRegistrationEventByID(ctx context.Context, id uint) (data RegistrationEvent, err error)
	GetRegistrationEventUserByStatus(ctx context.Context, eventID uint, userID string, statuses []string) (data []RegistrationEvent, err error)
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, payload CreateEventPayload) error
	UpdateEvent(ctx context.Context, id uint, payload UpdateEventPayload) error
	GetEvents(ctx context.Context, filter EventFilter) (data []EventDTO, pagination Pagination, err error)
	CreateRegistrationEvent(ctx context.Context, payload RegisterEventPayload) (RegisterEventResponse, error)
	CreateEventPay(ctx context.Context, payload EventPayPayload) error
	GetEventByID(ctx context.Context, id uint) (resp EventDTO, err error)
	GetEventByIDAdmin(ctx context.Context, id uint) (resp EventAdminDTO, err error)
	DeleteEvent(ctx context.Context, id uint) (err error)
	RegistrationStatus(ctx context.Context, orderNo string) (resp RegisterStatusResponse, err error)
	ListRegistration(ctx context.Context, filter EventFilter) (resp []RegistrationEvent, pagination Pagination, err error)
	ListEventPay(ctx context.Context, filter EventFilter) (data []EventPay, pagination Pagination, err error)
	PayProcess(ctx context.Context, payload PayProcessPayload) error
	ListRegistrationByEvent(ctx context.Context, eventID uint, filterPagination FilterPagination) (data []EventRegistrationDTO, pagination Pagination, err error)
	UpdateRegistrationStatus(ctx context.Context, registrationID uint, payload UpdateRegistrationStatusRequest) error
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
	UpdateEvent(w http.ResponseWriter, r *http.Request)
	ListRegistrationByEvent(w http.ResponseWriter, r *http.Request)
	UpdateRegistrationStatus(w http.ResponseWriter, r *http.Request)
}

type Event struct {
	ID                   uint                `json:"id" gorm:"primarykey"`
	Title                string              `json:"title" `
	Description          string              `json:"description"`
	Slug                 string              `json:"slug"`
	AuthorID             int                 `json:"author_id"`
	Author               User                `gorm:"foreignKey:AuthorID;references:ID"` // Ensure foreign key is correctly referenced
	Image                string              `json:"image"`
	Date                 null.Time           `json:"date"`
	Type                 constants.EventType `json:"type"`
	Location             string              `json:"location"`
	Duration             string              `json:"duration"`
	Capacity             int                 `json:"capacity"`
	Status               string              `json:"status"`
	Tags                 []EventTag          `json:"tags" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	Speakers             []EventSpeaker      `json:"speakers" gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE;"`
	RegistrationLink     string              `json:"registration_link"`
	Price                float64             `json:"price"` // 0 == free
	CreatedBy            int                 `json:"-"`
	UpdatedBy            int                 `json:"-"`
	DeletedBy            int                 `json:"-"`
	ReservationStartDate null.Time           `json:"reservation_start_date"`
	ReservationEndDate   null.Time           `json:"reservation_end_date"`
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            null.Time           `json:"-"`
	DeletedAt            null.Time           `json:"-"`
	AdditionalLink       string              `json:"additional_link"`
	SessionType          string              `json:"session_type"`
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
	Title                string              `json:"title" validate:"required"`
	Description          string              `json:"description" validate:"required"`
	Author               string              `json:"author" validate:"required"`
	FileName             string              `json:"file_name" validate:"required"`
	Slug                 string              `json:"slug" validate:"required"`
	Date                 null.Time           `json:"date" validate:"required"`
	Type                 constants.EventType `json:"type" validate:"required"`
	Location             string              `json:"location" validate:"required"`
	Duration             string              `json:"duration" validate:"required"`
	Status               string              `json:"status" validate:"required"`
	Capacity             int                 `json:"capacity" validate:"required"`
	Price                float64             `json:"price"`
	RegistrationLink     string              `json:"registration_link"`
	Tags                 []string            `json:"tags"`
	Speakers             []string            `json:"speakers"`
	ReservationStartDate null.Time           `json:"reservation_start_date"`
	ReservationEndDate   null.Time           `json:"reservation_end_date"`
	AdditionalLink       string              `json:"additional_link"`
	SessionType          string              `json:"session_type" validate:"required"`
}

type UpdateEventPayload struct {
	Title                string              `json:"title" validate:"required"`
	Description          string              `json:"description" validate:"required"`
	Author               string              `json:"author" validate:"required"`
	FileName             string              `json:"file_name" validate:"required"`
	Slug                 string              `json:"slug" validate:"required"`
	Date                 null.Time           `json:"date" validate:"required"`
	Type                 constants.EventType `json:"type" validate:"required"`
	Location             string              `json:"location" validate:"required"`
	Duration             string              `json:"duration" validate:"required"`
	Status               string              `json:"status" validate:"required"`
	Capacity             int                 `json:"capacity" validate:"required"`
	Price                float64             `json:"price"`
	RegistrationLink     string              `json:"registration_link"`
	Tags                 []string            `json:"tags"`
	Speakers             []string            `json:"speakers"`
	ReservationStartDate null.Time           `json:"reservation_start_date"`
	ReservationEndDate   null.Time           `json:"reservation_end_date"`
	AdditionalLink       string              `json:"additional_link"`
	SessionType          string              `json:"session_type" validate:"required"`
}

type EventDTO struct {
	ID               int                 `json:"id"`
	Title            string              `json:"title"`
	Description      string              `json:"description"`
	Author           string              `json:"author"`
	ImageEvent       string              `json:"image_event"`
	DateEvent        null.Time           `json:"date_event"`
	Type             constants.EventType `json:"type"`
	Price            float64             `json:"price"`
	Location         string              `json:"location"`
	Duration         string              `json:"duration"`
	Capacity         int                 `json:"capacity"`
	RegistrationLink string              `json:"registration_link"`
	SessionType      string              `json:"session_type"`
	Status           string              `json:"status"`
}

type EventAdminDTO struct {
	ID                   int                 `json:"id"`
	Title                string              `json:"title"`
	Description          string              `json:"description"`
	Author               string              `json:"author"`
	FileName             string              `json:"file_name"`
	Slug                 string              `json:"slug"`
	Date                 null.Time           `json:"date"`
	Type                 constants.EventType `json:"type"`
	Location             string              `json:"location"`
	Duration             string              `json:"duration"`
	Status               string              `json:"status"`
	Capacity             int                 `json:"capacity"`
	Price                float64             `json:"price"`
	RegistrationLink     string              `json:"registration_link"`
	Tags                 []string            `json:"tags"`
	Speakers             []string            `json:"speakers"`
	ReservationStartDate null.Time           `json:"reservation_start_date"`
	ReservationEndDate   null.Time           `json:"reservation_end_date"`
	// AdditionalLink       string              `json:"additional_link"`
	SessionType string `json:"session_type"`
}

type UpdateEvenPayload struct {
	Title            string              `json:"title"`
	Description      string              `json:"description"`
	Author           string              `json:"author"`
	ImageEvent       string              `json:"image_event"`
	DateEvent        null.Time           `json:"date_event"`
	Type             constants.EventType `json:"type"`
	Location         string              `json:"location"`
	Duration         string              `json:"duration"`
	Capacity         int                 `json:"capacity"`
	RegistrationLink string              `json:"registration_link"`
}

type EventFilter struct {
	ID        uint
	Title     string
	Type      constants.EventType
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

type EventPayPayload struct {
	OrderNo           string  `json:"order_no"`
	ImageProofPayment string  `json:"image_proof_payment"`
	NetAmount         float64 `json:"net_amount"`
}

type PayProcessPayload struct {
	OrderNo string `json:"order_no"`
	Status  string `json:"status"`
	Note    string `json:"note"`
}

func (e Event) ToDTO() EventDTO {
	// processing data
	id := int(e.ID)

	return EventDTO{
		ID:               id,
		Title:            e.Title,
		Description:      e.Description,
		Author:           e.Author.Username,
		ImageEvent:       e.Image,
		DateEvent:        e.Date,
		Type:             e.Type,
		Status:           e.Status,
		Location:         e.Location,
		Duration:         e.Duration,
		Capacity:         e.Capacity,
		Price:            e.Price,
		RegistrationLink: e.RegistrationLink,
		SessionType:      e.SessionType,
	}
}

func (e Event) ToAdminDTO() EventAdminDTO {
	id := int(e.ID)

	tags := make([]string, len(e.Tags))
	for i, tag := range e.Tags {
		tags[i] = tag.Tag
	}

	speakers := make([]string, len(e.Speakers))
	for i, speaker := range e.Speakers {
		speakers[i] = speaker.Name
	}

	return EventAdminDTO{
		ID:                   id,
		Title:                e.Title,
		Description:          e.Description,
		Author:               e.Author.Username,
		FileName:             e.Image,
		Slug:                 e.Slug,
		Date:                 e.Date,
		Type:                 e.Type,
		Location:             e.Location,
		Duration:             e.Duration,
		Status:               e.Status,
		Capacity:             e.Capacity,
		Price:                e.Price,
		RegistrationLink:     e.RegistrationLink,
		Tags:                 tags,
		Speakers:             speakers,
		ReservationStartDate: e.ReservationStartDate,
		ReservationEndDate:   e.ReservationEndDate,
		// AdditionalLink:       e.AdditionalLink,
		SessionType: e.SessionType,
	}
}
