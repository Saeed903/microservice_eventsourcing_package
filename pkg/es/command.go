package es

// Command commands interface for event sourcing
type Command interface {
	GetAggregateID() string
}

type BaseCommand struct {
	AggregateID string `json:"aggregateID" validate:"required,gte=0"`
}

func NewBaseCommand(aggregateId string) BaseCommand {
	return BaseCommand{AggregateID: aggregateId}
}

func (bc *BaseCommand) GetAggregateID() string {
	return bc.AggregateID
}
