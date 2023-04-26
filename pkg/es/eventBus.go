package es

import "context"

// EventBus ProcessEvents method publish events to the app specific message broker.
type EventBus interface {
	ProcessEvents(ctx context.Context, events []Event) error
}
