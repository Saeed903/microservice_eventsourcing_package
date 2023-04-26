package kafka

// Config kafka config
type Config struct {
	Brokers    []string `mapstructure:"brokers" validate:"required"`
	GroupId    string   `mapstructure:"groupID" validate:"required,gte=0"`
	InitTopics bool     `mapstructure:"initTopics"`
}

// TopicConfig kafka topic config
type TopicConfig struct {
	TopicName         string `mapstructure:"topicName" validate:"required"`
	Partitions        int    `mapstructure:"partitions" validate:"required,gte=0"`
	ReplicationFactor int    `mapstructure:"replicationFactor" validate:"required,gte=0"`
}
