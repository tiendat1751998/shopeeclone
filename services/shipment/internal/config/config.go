package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	AppName  string
	AppEnv   string
	LogLevel string
	HTTPPort int
	GRPCPort int
	MySQL    MySQLConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	JWT      JWTConfig
	Shipment ShipmentConfig
	Idempotency IdempotencyConfig
	OpenTelemetry OTELConfig
}

type MySQLConfig struct {
	Host         string
	Port         int
	User         string
	Password     string
	Database     string
	MaxOpenConns int
	MaxIdleConns int
	MaxLifetime  time.Duration
	Timeout      time.Duration
}

func (c MySQLConfig) DSN() string {
	return c.User + ":" + c.Password + "@tcp(" + c.Host + ":" + strconv.Itoa(c.Port) + ")/" + c.Database + "?charset=utf8mb4&parseTime=true&loc=UTC&timeout=" + c.Timeout.String()
}

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxRetries   int
}

type KafkaConfig struct {
	Brokers       []string
	TopicPrefix   string
	ConsumerGroup string
	DLQTopic      string
}

type JWTConfig struct {
	AccessSecret string
	AccessTTL    time.Duration
	Issuer       string
	Audience     string
}

type ShipmentConfig struct {
	DefaultCarrier    string
	IdempotencyTTL    time.Duration
	WebhookSecret     string
	TrackingSyncInterval time.Duration
}

type IdempotencyConfig struct {
	TTL time.Duration
}

type OTELConfig struct {
	Endpoint    string
	ServiceName string
	TraceRatio  float64
}

func Load() *Config {
	return &Config{
		AppName: getEnv("APP_NAME", "shopee-shipment"), AppEnv: getEnv("APP_ENV", "development"),
		LogLevel: getEnv("LOG_LEVEL", "info"), HTTPPort: getEnvInt("SHIPMENT_HTTP_PORT", 8085),
		GRPCPort: getEnvInt("SHIPMENT_GRPC_PORT", 9095),
		MySQL: MySQLConfig{
			Host: getEnv("MYSQL_HOST", "localhost"), Port: getEnvInt("MYSQL_PORT", 3306),
			User: getEnv("MYSQL_USER", "shopee"), Password: getEnv("MYSQL_PASSWORD", "shopee_dev"),
			Database: getEnv("MYSQL_DATABASE", "shopee_shipments"), MaxOpenConns: 25, MaxIdleConns: 10,
			MaxLifetime: 5 * time.Minute, Timeout: 5 * time.Second,
		},
		Redis: RedisConfig{
			Addr: getEnv("REDIS_ADDR", "localhost:6379"), Password: getEnv("REDIS_PASSWORD", ""),
			DB: getEnvInt("REDIS_DB", 5), PoolSize: 100, MinIdleConns: 20,
			DialTimeout: 5 * time.Second, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second, MaxRetries: 3,
		},
		Kafka: KafkaConfig{
			Brokers: getEnvSlice("KAFKA_BROKERS", ","), TopicPrefix: getEnv("KAFKA_TOPIC_PREFIX", "shopee.shipments"),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "shopee-shipment-service"), DLQTopic: "shopee.shipments.dlq",
		},
		JWT: JWTConfig{
			AccessSecret: getEnv("JWT_ACCESS_SECRET", "change-me-in-production"),
			AccessTTL: getEnvDuration("JWT_ACCESS_TTL", 15*time.Minute), Issuer: "shopee-auth", Audience: "shopee-clone",
		},
		Shipment: ShipmentConfig{
			DefaultCarrier: getEnv("DEFAULT_CARRIER", "ninja_van"),
			IdempotencyTTL: getEnvDuration("SHIPMENT_IDEMPOTENCY_TTL", 24*time.Hour),
			WebhookSecret: getEnv("WEBHOOK_SECRET", "whsec-change-me"),
			TrackingSyncInterval: getEnvDuration("TRACKING_SYNC_INTERVAL", 5*time.Minute),
		},
		Idempotency: IdempotencyConfig{TTL: getEnvDuration("IDEMPOTENCY_TTL", 24*time.Hour)},
		OpenTelemetry: OTELConfig{
			Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "http://localhost:4318"),
			ServiceName: getEnv("OTEL_SERVICE_NAME", "shopee-shipment"),
			TraceRatio: getEnvFloat("OTEL_TRACES_SAMPLER_ARG", 0.1),
		},
	}
}

func getEnv(key, fallback string) string { if v := os.Getenv(key); v != "" { return v }; return fallback }
func getEnvInt(key string, fallback int) int { if v := os.Getenv(key); v != "" { if i, err := strconv.Atoi(v); err == nil { return i } }; return fallback }
func getEnvDuration(key string, fallback time.Duration) time.Duration { if v := os.Getenv(key); v != "" { if d, err := time.ParseDuration(v); err == nil { return d } }; return fallback }
func getEnvFloat(key string, fallback float64) float64 { if v := os.Getenv(key); v != "" { if f, err := strconv.ParseFloat(v, 64); err == nil { return f } }; return fallback }
func getEnvSlice(key, sep string) []string { v := getEnv(key, ""); if v == "" { return []string{"localhost:9092"} }; return strings.Split(v, sep) }
