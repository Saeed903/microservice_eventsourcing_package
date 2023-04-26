package kafka

import (
	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
)

// NewWriter create new configured kafka writer
func NewWriter(brokers []string, errLogger kafka.Logger) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  writerMaxAttempts,
		ErrorLogger:  errLogger,
		Compression:  compress.Snappy,
		ReadTimeout:  writeReadTimeout,
		WriteTimeout: writerWriteTimeout,
		BatchTimeout: batchTimeout,
		BatchSize:    batchSize,
		Async:        false,
	}
}

// NewAsyncWriter Create new configured kafka async writer
func NewAsyncWriter(brokers []string, errLogger kafka.Logger, log logger.Logger) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  maxAttempts,
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  errLogger,
		Compression:  compress.Snappy,
		ReadTimeout:  writeReadTimeout,
		WriteTimeout: writerWriteTimeout,
		Async:        true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Errorf("(kafka.AsyncWriter Error) topic: %s, partition: %v, offset: %v err: %v", messages[0].Topic,
					messages[0].Partition, messages[0].Offset, err)
				return
			}
		},
	}
}

type AsyncWriterCallback func(messages []kafka.Message) error

// NewAsyncWriterWithCallback create new configured kafka async writer with callback function
func NewAsyncWriterWithCallback(brokers []string, errLogger kafka.Logger, log logger.Logger, cb AsyncWriterCallback) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		MaxAttempts:  maxAttempts,
		RequiredAcks: kafka.RequireAll,
		ErrorLogger:  errLogger,
		Compression:  compress.Snappy,
		ReadTimeout:  writeReadTimeout,
		WriteTimeout: writerWriteTimeout,
		Async:        true,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Errorf("(kafka.AsyncWriter Error) topic: %s, partition: %v, offset: %v", messages[0].Topic, messages[0].Partition, messages[0].Offset)
				if err := cb(messages); err != nil {
					log.Errorf("(kafka.AsyncWriter Callback Error) err: %v", err)
					return
				}
				return
			}
		},
	}
}

// NewRequireNoneWriter create new configured kafka writer
func NewRequireNoneWriter(brokers []string, errLogger kafka.Logger, log logger.Logger) *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireAll,
		MaxAttempts:  maxAttempts,
		ErrorLogger:  errLogger,
		Compression:  compress.Snappy,
		ReadTimeout:  writerRequireNoneReadTimeout,
		WriteTimeout: writerRequireNonWriterTimeout,
		Async:        false,
		Completion: func(messages []kafka.Message, err error) {
			if err != nil {
				log.Errorf("(kafka.Writer Error) topic: %s, partition: %v, offset: %v", messages[0].Topic, messages[0].Partition, messages[0].Offset)
				return
			}
		},
	}
}
