package game

import "github.com/google/uuid"

type SerializedEvent struct {
	id        uuid.UUID
	version   int
	eventName string
	event     []byte
}

type EventRepository interface {
	StoreEvents(id uuid.UUID, version int, events []Event) error
	RetrieveEvents(id uuid.UUID) ([]Event, error)
}
