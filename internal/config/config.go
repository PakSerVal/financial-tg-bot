package config

import (
	"flag"
)

const (
	defaultPgHost        = "localhost"
	defaultPgPort        = 5432
	defaultPgUser        = "postgres"
	defaultPgPass        = "postgres"
	defaultPgDatabase    = "postgres"
	defaultSslMode       = "disable"
	defaultRedisHost     = "localhost"
	defaultRedisPort     = 6379
	defaultRedisPassword = ""
	defaultRedisDb       = 0
)

type Config struct {
	token       string
	useInmemory bool
	dev         bool
	serviceName string
	httpPort    int64
	dbConn      ConnConfig
	redisConn   RedisConfig
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

func New() (*Config, error) {
	c := &Config{}

	flag.StringVar(&c.token, "token", "", "bot token")
	flag.BoolVar(&c.useInmemory, "useInmemory", false, "inmemory usage")
	flag.BoolVar(&c.dev, "dev", false, "Development mode")
	flag.StringVar(&c.serviceName, "serviceName", "telegram bot", "Service name")
	flag.Int64Var(&c.httpPort, "httpPort", 8080, "Http port")
	parseDbConn(&c.dbConn)
	parseRedisConn(&c.redisConn)

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

func (c *Config) RedisConfig() RedisConfig {
	return c.redisConn
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
