package constants

const (
	Open   = "open"
	Soon   = "soon"
	Closed = "closed"
)

type EventType string

const (
	EventTypeConference EventType = "Conference"
	EventTypeTechTalk   EventType = "Tech Talk"
	EventTypeNgobar     EventType = "Ngobar"

	PENDING   = "PENDING"
	SUCCESS   = "SUCCESS"
	REJECTED  = "REJECTED"
	CANCELLED = "CANCELLED"
	EXPIRED   = "EXPIRED"
)

func GetValidEventTypes() []EventType {
	return []EventType{
		EventTypeConference,
		EventTypeTechTalk,
		EventTypeNgobar,
	}
}

func IsValidEventType(eventType EventType) bool {
	validTypes := GetValidEventTypes()
	for _, validType := range validTypes {
		if validType == eventType {
			return true
		}
	}
	return false
}
