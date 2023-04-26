package constants

const (
	GrpcPort         = "GRPC_PORT"
	HttpPort         = "HTTP_PORT"
	ConfigPath       = "CONFIG_PATH"
	KafkaBrokers     = "KAFKA_BROKERS"
	JaegerHostPort   = "JAEGER_HOST"
	RedisAddr        = "REDIS_HOST"
	MongoDbURI       = "MONGO_DB_URI"
	PostgresqlHost   = "POSTGRE_HOST"
	PostgresqlPort   = "POSTGRE_PORT"
	ElasticUrl       = "ELASTIC_URL"
	TCP              = "TCP"
	MIGRATION_DB_URL = "MIGRATION_DB_URL"

	ReaderServicePort = "READER_SERVICE_PORT"

	Yaml          = "yaml"
	Redis         = "redis"
	Kafka         = "kafka"
	PostgreSql    = "postgreSql"
	MongoDb       = "mongo"
	ElasticSearch = "elasticSearch"

	GRPC     = "GRPC"
	SIZE     = "SIZE"
	URI      = "URI"
	STATUS   = "STATUS"
	HTTP     = "HTTP"
	ERROR    = "ERROR"
	METHOD   = "METHOD"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic     = "topic"
	Partition = "partition"
	Message   = "message"
	WorkerID  = "workerID"
	Headers   = "headers"
	Offset    = "offset"
	Time      = "time"

	Page   = "page"
	Size   = "size"
	Search = "search"
	ID     = "id"

	EsAll = "$all"

	Validate        = "validate"
	FieldValidation = "field Validation"
	RequiredHeaders = "required header"
	Base64          = "base64"
	Unmarshal       = "unmarshal"
	Uuid            = "uuid"
	Cookie          = "cookie"
	Token           = "token"
	Bcrypt          = "bcrypt"
	SQLState        = "sqlstate"

	MongoAggregateID = "aggregateID"

	MongoProjection   = "(MongoDB Projection)"
	ElasticProjection = "(Elastic Projection)"

	EventType   = "(EventType)"
	AggregateID = "(AggregateID)"
	Version     = "(Version)"
	TimeStamp   = "(TimeStamp)"
	Metadata    = "(Metadata)"

	MessageSize = "MessageSize"

	ManagementRiskIndex = "aggregateID"
	ManagementRiskId    = "ManagementRiskId"

	KafkaHeaders = "KafkaHeaders"

	Tcp = "tcp"
)
