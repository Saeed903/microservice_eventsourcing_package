package es

import "github.com/pkg/errors"

var (
	ErrAlreadyExists       = errors.New("Already exists")
	ErrAggregateNotFound   = errors.New("Aggregate not found")
	ErrInvalidEventType    = errors.New("Invalid event type")
	ErrInvalidCommandType  = errors.New("Invalid command type")
	ErrInvalidAggregate    = errors.New("Invalid aggregate")
	ErrInvalidAggregateID  = errors.New("Invalid aggregateid")
	ErrInvalidEventVersion = errors.New("Invalid event version")
)
