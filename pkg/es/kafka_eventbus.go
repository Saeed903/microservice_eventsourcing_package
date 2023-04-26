package es

import (
	"context"
	"fmt"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/es/serializer"
	kafkaClient "github.com/saeed903/microservice_eventsourcing_package/pkg/kafka"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/tracing"
	"github.com/segmentio/kafka-go"
)

// KafkaEventBusConfig kafka eventbus config
type KafkaEventBusConfig struct {
	Topic             string `mapstructure:"topic" validae:"required"`
	TopicPerfix       string `mapstructure:"topicPerfic" validate:"required"`
	Partitions        int    `mapstructure:"partitions" validate:"required,gte=0"`
	ReplicationFactor int    `mapstructure:"replicationFactor" validate:"required,gte=0"`
	Headers           []kafka.Header
}

type KafkaEventsBus struct {
	producer kafkaClient.Producer
	cfg      KafkaEventBusConfig
}

// NewKafkaEventsBus kafkaEventsBus constructor.
func NewKafkaEventsBus(producer kafkaClient.Producer, cfg KafkaEventBusConfig) *KafkaEventsBus {
	return &KafkaEventsBus{
		producer: producer,
		cfg:      cfg,
	}
}

// ProcessEvents serialize to json and publish es.Event's to the kafka bus.
func (e *KafkaEventsBus) ProcessEvents(ctx context.Context, events []Event) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "KafkaEventsBus.ProcessEvents")
	defer span.Finish()

	eventsBytes, err := serializer.Marshal(events)
	if err != nil {
		return tracing.TraceWithErr(span, errors.Wrap(err, "serializer.Marshal"))
	}

	return e.producer.PublicMessage(ctx, kafka.Message{
		Topic:   GetTopicName(e.cfg.TopicPerfix, string(events[0].GetAggregateType())),
		Value:   eventsBytes,
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
		Time:    time.Now().UTC(),
	})
}

func GetTopicName(eventStorePerfix, aggregateType string) string {
	return fmt.Sprintf("%s_%s", eventStorePerfix, aggregateType)
}

func GetkafkaAggregateTypeTopic(cfg KafkaEventBusConfig, aggregateType string) kafka.TopicConfig {
	return kafka.TopicConfig{
		Topic:             GetTopicName(cfg.TopicPerfix, aggregateType),
		NumPartitions:     cfg.Partitions,
		ReplicationFactor: cfg.ReplicationFactor,
	}
}
