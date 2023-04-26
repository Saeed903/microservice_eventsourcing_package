package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/constants"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/elastic"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/es"
	kafkaClient "github.com/saeed903/microservice_eventsourcing_package/pkg/kafka"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/logger"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/migrations"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/mongodb"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/postgres"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/probes"
	"github.com/saeed903/microservice_eventsourcing_package/pkg/tracing"
	"github.com/spf13/viper"
)

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "", "ManagementRisk microservice pah")
}

type Config struct {
	ServiceName          string                 `mapstructure:"serviceName"`
	Logger               logger.LogConfig       `mapstructure:"logger"`
	GRPC                 GRPC                   `mapstructure:"grpc"`
	Postgresql           postgres.Config        `mapstructure:"postgres"`
	Timeouts             Timeouts               `mapstructure:"timeouts" validate:"required"`
	EventSourcingConfig  es.Config              `mapstructure:"eventSourcingConfig" validate:"required"`
	Kafka                *kafkaClient.Config    `mapstructure:"kafka" validate:"required"`
	KafkaTopics          KafkaTopics            `mapstructure:"kafkaTopics" validate:"required"`
	Mongo                *mongodb.Config        `mapstructre:"mongo" validate:"required"`
	MongoCollections     MongoCollections       `mapstructre:"mongoCollections" validate:"required"`
	KafkaPublisherConfig es.KafkaEventBusConfig `mapstructre:"kafkaPublisherConfig" validate:"required"`
	Jaeger               *tracing.Config        `mapstructure:"jaeger"`
	ElasticIndexes       ElasticIndexes         `mapstructure:"elasticIndexes" validate:"required"`
	Projections          Projections            `mapstructure:"projections"`
	Http                 Http                   `mapstructure:"http"`
	Probes               probes.Config          `mapstructure:"probes"`
	ElasticSearch        elastic.Config         `mapstructure:"elasticSearch" validate:"required"`
	MigrationConfig      migrations.Config      `mapstructre:"migrationConfig" validate:"required"`
}

type GRPC struct {
	Port        string `mapstructure:"port"`
	Development bool   `mapstructure:"development"`
}

type Timeouts struct {
	PostgresInitMiliseconds int  `mapstructure:"postgresInitMiliseconds" validate:"required"`
	PostgresInitRetryCount  uint `mapstructure:"postgresInitRetryCount" validate:"required"`
}

type KafkaTopics struct {
	EventCreated                     kafkaClient.TopicConfig `mapstructure:"eventCreated" validate:"required"`
	ManagementRiskAggregateTypeTopic kafkaClient.TopicConfig `mapstructure:"managementRiskAggregateTypeTopic" validate:"required"`
}

type MongoCollections struct {
	ManagementRisk string `mapstructure:"managementRisk" validate:"required"`
}

type ElasticIndexes struct {
	ManagementRisk string `mapstructure:"managementRisk" validate:"required"`
}

type Projections struct {
	MongoGroup                  string `mapstructre:"mongoGroup" validate:"required"`
	MongoSubscriptionPoolSize   int    `mapstructure:"mongoSubscriptionPoolSize" validate:"required,gte=0"`
	ElasticGroup                string `mapstructure:"elasticGroup" validate:"required"`
	ElasticSubscriptionPoolSize int    `mapstructure:"elasticSubscriptionPoolSize" validate:"required,gte=0"`
}

type Http struct {
	Port                string   `mapstructure:"port" validate:"required"`
	Development         bool     `mapstructure:"development"`
	BasePath            string   `mapstructure:"basePath" validate:"required"`
	ManagementRiskPath  string   `mapstructure:"managementRiskPath" validate:"require"`
	DebugErrorsResponse bool     `mapstructure:"debugErrorsResponse"`
	IgnorLogUrls        []string `mapstructure:"ignorLogUrls"`
}

func InitConfig() (*Config, error) {
	if configPath == "" {
		configPathFromEnv := os.Getenv(constants.ConfigPath)
		if configPathFromEnv != "" {
			configPath = configPathFromEnv
		} else {
			getwd, err := os.Getwd()
			if err != nil {
				return nil, errors.Wrap(err, "os.Getwd")
			}
			configPath = fmt.Sprintf("%s/config/config.yaml", getwd)
		}
	}

	cfg := &Config{}

	viper.SetConfigType(constants.Yaml)
	viper.SetConfigFile(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Wrap(err, "viper.ReadInConfig")
	}

	if err := viper.Unmarshal(cfg); err != nil {
		return nil, errors.Wrap(err, "viper.Unmarshal")
	}

	grpcPort := os.Getenv(constants.GrpcPort)
	if grpcPort != "" {
		cfg.GRPC.Port = grpcPort
	}

	mongoURI := os.Getenv(constants.MongoDbURI)
	if mongoURI != "" {
		cfg.Mongo.URI = mongoURI
	}

	jaegerAddr := os.Getenv(constants.JaegerHostPort)
	if jaegerAddr != "" {
		cfg.Jaeger.HostPort = jaegerAddr
	}

	elasticUrl := os.Getenv(constants.ElasticUrl)
	if elasticUrl != "" {
		cfg.ElasticSearch.Addresses = []string{elasticUrl}
	}

	postgresPort := os.Getenv(constants.PostgresqlPort)
	if postgresPort != "" {
		cfg.Postgresql.Port = postgresPort
	}

	dbUrl := os.Getenv(constants.MIGRATION_DB_URL)
	if dbUrl != "" {
		cfg.MigrationConfig.DbURL = dbUrl
	}

	kafkaBrokers := os.Getenv(constants.KafkaBrokers)
	if kafkaBrokers != "" {
		cfg.Kafka.Brokers = []string{kafkaBrokers}
	}

	return cfg, nil

}
