package es

import "fmt"

const (
	changeEventCap = 10
	startVersion   = 0
)

type When interface {
	When(event any) error
}

type when func(event any) error

// Apply process Aggregate Event
type Apply interface {
	Apply(event any) error
}

// Load create Aggregate state from Event's
type Load interface {
	Load(events []any) error
}

// RaiseEvent process applied Aggregate Event from event store
type RaiseEvent interface {
	RaiseEvent(event any) error
}

type Aggregate interface {
	When
	AggregateRoot
	RaiseEvent
	Load
	Apply
}

type AggregateRoot interface {
	GetID() string
	SetID(id string) *AggregateBase
	GetType() AggregateType
	SetType(aggregateType AggregateType)
	GetChanges() []any
	ClearChanges()
	GetVersion() uint64
	ToSnapshot()
	String() string
}

// AggregateType type of the aggregate
type AggregateType string

// AggregateBase base aggregate contains all main necessary fields
type AggregateBase struct {
	ID      string
	Version uint64
	Changes []any
	Type    AggregateType
	when    when
}

// NewAggregateBase constructor, contains all main fields and methods,
// main aggregate must realize When interface and pass as arguments to constructor
// Example of recommended aggregate constructor method:
//
// func NewOrderAggregate() *OrderAggregate {
// 		orderAggregate := &OrderAggregate{
//				Order : Models.NewOrder(),
//		}
// 		base := es.NewAggregateBase(orderAggregate.When)
//		base.SetType(OrderAggregateType)
//		orderAggregate.AggregateBase = base
//		return orderAggregate
// }
func NewAggregateBase(when when) *AggregateBase {
	if when == nil {
		return nil
	}

	return &AggregateBase{
		Version: startVersion,
		Changes: make([]any, 0, changeEventCap),
		when:    when,
	}
}

// SetID set AggregateBase ID
func (a *AggregateBase) SetID(id string) *AggregateBase {
	a.ID = id
	return a
}

// GetID get AggregateBase ID
func (a *AggregateBase) GetID() string {
	return a.ID
}

// SetType set AggreagteBase AggregateType
func (a *AggregateBase) SetType(aggregateType AggregateType) {
	a.Type = aggregateType
}

// GetType get AggregateBase AggregateType
func (a *AggregateBase) GetType() AggregateType {
	return a.Type
}

// GetVersion get AggregateBase version
func (a *AggregateBase) GetVersion() uint64 {
	return a.Version
}

// ClearChanges clear AggregateBase uncommitted Event's
func (a *AggregateBase) ClearChanges() {
	a.Changes = make([]any, 0, changeEventCap)
}

// GetChanges get AggregateBase uncommitted Event's
func (a *AggregateBase) GetChanges() []any {
	return a.Changes
}

// Load add existing events from event store to aggregate using When interface method
func (a *AggregateBase) Load(events []any) error {
	for _, evt := range events {
		if err := a.when(evt); err != nil {
			return err
		}

		a.Version++
	}

	return nil
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event any) error {
	if err := a.when(event); err != nil {
		return err
	}

	a.Version++
	a.Changes = append(a.Changes, event)
	return nil
}

// RaiseEvent push event to aggregate applied events using When method, used for load directly from eventstore
func (a *AggregateBase) RaiseEvent(event any) error {
	if err := a.when(event); err != nil {
		return err
	}

	a.Version++
	return nil
}

// ToSnapshot prepare AggregateBase for saving Snapshot.
func (a *AggregateBase) ToSnapshot() {
	a.ClearChanges()
}

func (a *AggregateBase) String() string {
	return fmt.Sprintf("(Aggregate) AggregateID: %s, Type: %s, Version: %d, Changes: %d",
		a.GetID(),
		string(a.GetType()),
		a.GetVersion(),
		len(a.GetChanges()),
	)
}
