package kafka

import (
	"context"
	"sync"

	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"
	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/compress"
	"golang.org/x/sync/errgroup"
)

// MessageProcessor processor must implement kafka.Worker func method interface.
type MessageProcessor interface {
	ProcessMessages(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)
	ProcessMessagesWithErrGroup(ctx context.Context, r *kafka.Reader, workerID int)
}

// Worker kafka consumer worker fetch and process messages form  reader
type Worker func(ctx context.Context, r *kafka.Reader, wg *sync.WaitGroup, workerID int)

// WorkerErrGroup kafka consumer worker fetch and process messages from reader
type WorkerErrGroup func(ctx context.Context, r *kafka.Reader, workerID int) error

type ConsumerGroup interface {
	ConsumeTopic(ctx context.Context, groupTopics []string, poolSize int, worker Worker)
	ConsumeTopicWithErrGroup(ctx context.Context, groupTopics []string, poolSize int, worker WorkerErrGroup) error
	GetNewKafkaReader(kafkaURL []string, groupTopics []string, groupID string) *kafka.Reader
	GetNewKafkaWriter() *kafka.Writer
}

type consumerGroup struct {
	Brokers []string
	GroupID string
	log     logger.Logger
}

// NewConsumerGroup kafka consumer group constructor
func NewConsumerGroup(brokers []string, groupID string, log logger.Logger) *consumerGroup {
	return &consumerGroup{
		log:     log,
		Brokers: brokers,
		GroupID: groupID,
	}
}

// GetNewKafkaReader create new kafka reader
func (c *consumerGroup) GetNewKafkaReader(kafkaURL []string, groupTopics []string, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:                kafkaURL,
		GroupID:                groupID,
		GroupTopics:            groupTopics,
		MinBytes:               minBytes,
		MaxBytes:               maxBytes,
		QueueCapacity:          queueCapacity,
		HeartbeatInterval:      heartbeatInterval,
		CommitInterval:         commitInterval,
		PartitionWatchInterval: maxAttempts,
		MaxWait:                maxWait,
		Dialer:                 &kafka.Dialer{Timeout: dialTimeout},
	})
}

// GetNewKafkaWriter create new kafka producer
func (c *consumerGroup) GetNewKafkaWriter() *kafka.Writer {
	return &kafka.Writer{
		Addr:         kafka.TCP(c.Brokers...),
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: writerRequiredAcks,
		MaxAttempts:  writerMaxAttempts,
		Compression:  compress.Snappy,
		ReadTimeout:  writeReadTimeout,
		WriteTimeout: writerWriteTimeout,
	}
}

// ConsumeTopic start consumer group with given worker and pool size
func (c *consumerGroup) ConsumeTopic(ctx context.Context, groupTopic []string, poolSize int, worker Worker) {
	r := c.GetNewKafkaReader(c.Brokers, groupTopic, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.log.Warnf("consumerGroup.r.Close: %v", err)
		}
	}()

	c.log.Infof("(Starting consumer groupID): GroupID %s, topics: %+v, poolSize: %v", c.GroupID, groupTopic, poolSize)

	wg := &sync.WaitGroup{}

	for i := 0; i < poolSize; i++ {
		wg.Add(1)
		go worker(ctx, r, wg, i)
	}
	wg.Wait()
}

// ConsumeTopicWithErrGroup start conusmer group with given worker and pool size
func (c *consumerGroup) ConsumeTopicWithErrGroup(ctx context.Context, groupTopics []string, poolSize int, worker WorkerErrGroup) error {
	r := c.GetNewKafkaReader(c.Brokers, groupTopics, c.GroupID)

	defer func() {
		if err := r.Close(); err != nil {
			c.log.Warnf("consumerGroup.r.Close: %v", err)
		}
	}()

	c.log.Infof("(Starting ConsumeTopicWithErrGroup) GroupID: %s, topics: %+v, poolSize: %d", c.GroupID, groupTopics, poolSize)

	g, ctx := errgroup.WithContext(ctx)
	for i := 0; i < poolSize; i++ {
		g.Go(c.runWorker(ctx, worker, r, i))
	}
	return g.Wait()

}

func (c *consumerGroup) runWorker(ctx context.Context, worker WorkerErrGroup, r *kafka.Reader, i int) func() error {
	return func() error {
		return worker(ctx, r, i)
	}
}