package config

import (
	"flag"
)

const (
	defaultPgHost              = "localhost"
	defaultPgPort              = 5432
	defaultPgUser              = "postgres"
	defaultPgPass              = "postgres"
	defaultPgDatabase          = "postgres"
	defaultSslMode             = "disable"
	defaultRedisHost           = "localhost"
	defaultRedisPort           = 6379
	defaultRedisPassword       = ""
	defaultRedisDb             = 0
	defaultKafkaBrokerList     = "localhost:9092"
	defaultKafkaVersion        = "2.5.0"
	defaultReturnSuccesses     = true
	defaultKafkaConsumerOffset = -2
	defaultKafkaConsumerTopics = "report"
	defaultKafkaGroupId        = "report-group"
	defaultKafkaAssignor       = "range"
)

type Config struct {
	token       string
	useInmemory bool
	dev         bool
	serviceName string
	httpPort    int64
	grpcHost    string
	grpcPort    int64
	dbConn      ConnConfig
	redisConn   RedisConfig
	kafkaConfig KafkaConfig
}

type ConnConfig struct {
	Host     string
	Port     int64
	User     string
	Password string
	DbName   string
	SslMode  string
}

type RedisConfig struct {
	Host     string
	Port     int64
	Password string
	Db       int64
}

type KafkaConfig struct {
	BrokerList string
	Version    string
	Producer   ProducerConfig
	Consumer   ConsumerConfig
}

type ProducerConfig struct {
	ReturnSuccesses bool
}

type ConsumerConfig struct {
	Offset   int64
	Topics   string
	GroupId  string
	Assignor string
}

func New() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.token, "token", "", "bot token")
	flag.BoolVar(&c.useInmemory, "useInmemory", false, "inmemory usage")
	flag.BoolVar(&c.dev, "dev", false, "Development mode")
	flag.StringVar(&c.serviceName, "serviceName", "telegram bot", "Service name")
	flag.Int64Var(&c.httpPort, "httpPort", 8080, "Http port")
	flag.StringVar(&c.grpcHost, "grpcHost", "localhost", "Grpc host")
	flag.Int64Var(&c.grpcPort, "grpcPort", 50051, "Grpc port")
	parseDbConn(&c.dbConn)
	parseRedisConn(&c.redisConn)
	parseKafkaConfig(&c.kafkaConfig)

	flag.Parse()

	return c, nil
}

func (c *Config) Token() string {
	return c.token
}

func (c *Config) DbConn() *ConnConfig {
	return &c.dbConn
}

func (c *Config) UseInmemory() bool {
	return c.useInmemory
}

func (c *Config) IsDev() bool {
	return c.dev
}

func (c *Config) ServiceName() string {
	return c.serviceName
}

func (c *Config) HttpPort() int64 {
	return c.httpPort
}

func (c *Config) GrpcHost() string {
	return c.grpcHost
}

func (c *Config) GrpcPort() int64 {
	return c.grpcPort
}

func (c *Config) RedisConfig() RedisConfig {
	return c.redisConn
}

func (c *Config) KafkaConfig() KafkaConfig {
	return c.kafkaConfig
}

func parseDbConn(c *ConnConfig) {
	flag.StringVar(&c.Host, "pgHost", defaultPgHost, "postgres host")
	flag.Int64Var(&c.Port, "pgPort", defaultPgPort, "postgres port")
	flag.StringVar(&c.User, "pgUser", defaultPgUser, "postgres user")
	flag.StringVar(&c.Password, "pgPass", defaultPgPass, "postgres password")
	flag.StringVar(&c.DbName, "pgDatabase", defaultPgDatabase, "postgres database name")
	flag.StringVar(&c.SslMode, "pgSslMode", defaultSslMode, "postgres ssl mode")
}

func parseRedisConn(r *RedisConfig) {
	flag.StringVar(&r.Host, "redisHost", defaultRedisHost, "cache host")
	flag.Int64Var(&r.Port, "redisPort", defaultRedisPort, "cache port")
	flag.StringVar(&r.Password, "redisPass", defaultRedisPassword, "cache password")
	flag.Int64Var(&r.Db, "redisDb", defaultRedisDb, "cache database")
}

func parseKafkaConfig(k *KafkaConfig) {
	flag.StringVar(&k.BrokerList, "kafkaBrokerList", defaultKafkaBrokerList, "kafka broker list")
	flag.StringVar(&k.Version, "kafkaVersion", defaultKafkaVersion, "kafka version")
	flag.BoolVar(&k.Producer.ReturnSuccesses, "kafkaProducerReturnSuccesses", defaultReturnSuccesses, "kafka producer return successes")
	flag.Int64Var(&k.Consumer.Offset, "kafkaConsumerOffset", defaultKafkaConsumerOffset, "kafka consumer offset")
	flag.StringVar(&k.Consumer.Topics, "kafkaConsumerTopics", defaultKafkaConsumerTopics, "kafka consumer topics")
	flag.StringVar(&k.Consumer.GroupId, "kafkaConsumerGroupId", defaultKafkaGroupId, "kafka consumer group id")
	flag.StringVar(&k.Consumer.Assignor, "kafkaConsumerAssignor", defaultKafkaAssignor, "kafka consumer assignor")
}
