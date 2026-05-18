package config
import ("os"; "strconv"; "time")
type Config struct { AppName string; AppEnv string; LogLevel string; HTTPPort int; GRPCPort int; MySQL MySQLConfig; Redis RedisConfig; Kafka KafkaConfig; OpenTelemetry OTELConfig }
type MySQLConfig struct { Host string; Port int; User string; Password string; Database string; MaxOpenConns int; MaxIdleConns int; MaxLifetime time.Duration; Timeout time.Duration }
func (c MySQLConfig) DSN() string { return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.Database + "?charset=utf8mb4&parseTime=true&loc=UTC&timeout=" + c.Timeout.String() }
type RedisConfig struct { Addr string; Password string; DB int; PoolSize int; MinIdleConns int; DialTimeout time.Duration; ReadTimeout time.Duration; WriteTimeout time.Duration; MaxRetries int }
type KafkaConfig struct { Brokers []string }
type OTELConfig struct { Endpoint string; ServiceName string; TraceRatio float64 }
func Load() *Config {
	return &Config{AppName: getEnv("APP_NAME", "shopee-notification"), AppEnv: getEnv("APP_ENV", "development"), LogLevel: getEnv("LOG_LEVEL", "info"), HTTPPort: getEnvInt("NOTIFY_HTTP_PORT", 8080), GRPCPort: getEnvInt("NOTIFY_GRPC_PORT", 9090),
		MySQL: MySQLConfig{Host: getEnv("MYSQL_HOST", "localhost"), Port: getEnvInt("MYSQL_PORT", 3306), User: getEnv("MYSQL_USER", "shopee"), Password: getEnv("MYSQL_PASSWORD", "shopee_dev"), Database: getEnv("MYSQL_DATABASE", "shopee_notification"), MaxOpenConns: 25, MaxIdleConns: 10, MaxLifetime: 5 * time.Minute, Timeout: 5 * time.Second},
		Redis: RedisConfig{Addr: getEnv("REDIS_ADDR", "localhost:6379"), Password: getEnv("REDIS_PASSWORD", ""), DB: getEnvInt("REDIS_DB", 0), PoolSize: 100, MinIdleConns: 20, DialTimeout: 5 * time.Second, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second, MaxRetries: 3},
		Kafka: KafkaConfig{Brokers: []string{getEnv("KAFKA_BROKERS", "localhost:9092")}},
		OpenTelemetry: OTELConfig{Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"), ServiceName: "shopee-notification", TraceRatio: 0.1}}
}
func getEnv(k, f string) string { if v := os.Getenv(k); v != "" { return v }; return f }
func getEnvInt(k string, f int) int { if v := os.Getenv(k); v != "" { if i, e := strconv.Atoi(v); e == nil { return i } }; return f }
func getEnvDuration(k string, f time.Duration) time.Duration { if v := os.Getenv(k); v != "" { if d, e := time.ParseDuration(v); e == nil { return d } }; return f }
