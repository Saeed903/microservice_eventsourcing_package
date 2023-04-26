package kafka

import (
	"context"

	"github.com/opentracing/opentracing-go"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/tracing"
	"github.com/segmentio/kafka-go"
)

type Producer interface {
	PublicMessage(ctx context.Context, msgs ...kafka.Message) error
	Close() error
}

type producer struct {
	log     logger.Logger
	brokers []string
	w       *kafka.Writer
}

// NewProducer create new kafka producer
func NewProducer(log logger.Logger, brokers []string) *producer {
	return &producer{
		log:     log,
		brokers: brokers,
		w:       NewWriter(brokers, kafka.LoggerFunc(log.Errorf)),
	}
}

// NewAsyncProducer create new kafka producer
func NewAsyncProducer(log logger.Logger, brokers []string) *producer {
	return &producer{
		log:     log,
		brokers: brokers,
		w:       NewAsyncWriter(brokers, kafka.LoggerFunc(log.Errorf), log),
	}
}

// NewAsyncProducerWithCallback create new kafka producer with callback for delete invalid projection
func NewAsyncProducerWithCallback(log logger.Logger, brokers []string, cb AsyncWriterCallback) *producer {
	return &producer{
		log:     log,
		brokers: brokers,
		w:       NewAsyncWriterWithCallback(brokers, kafka.LoggerFunc(log.Errorf), log, cb),
	}
}

// NewRequireNoneProducer create new fire and forget kafka producer
func NewRequireNoneProducer(log logger.Logger, brokers []string) *producer {
	return &producer{
		log:     log,
		brokers: brokers,
		w:       NewRequireNoneWriter(brokers, kafka.LoggerFunc(log.Errorf), log),
	}
}

// PublishMessage create publish message to topic kafka
func (p *producer) PublicMessage(ctx context.Context, msgs ...kafka.Message) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "prodcuer.PublishMessage")
	defer span.Finish()

	if err := p.w.WriteMessages(ctx, msgs...); err != nil {
		tracing.TraceErr(span, err)
		return err
	}
	return nil
}

func (p *producer) Close() error {
	return p.w.Close()
}
