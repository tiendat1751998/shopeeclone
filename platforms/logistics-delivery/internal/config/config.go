package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName  string
	Env      string
	Version  string
	HTTPPort string
	GRPCPort string

	Postgres PostgresConfig
	Redis    RedisConfig
	Kafka    KafkaConfig
	OTEL     OTELConfig
}

type PostgresConfig struct {
	DSN             string
	MaxConns        int
	MinConns        int
	MaxConnLifetime time.Duration
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type KafkaConfig struct {
	Brokers         []string
	ShipmentTopic   string
	TrackingTopic   string
	CourierTopic    string
	FulfillmentTopic string
	ConsumerGroup   string
}

type OTELConfig struct {
	Endpoint string
	Insecure bool
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}

func Load() *Config {
	return &Config{
		AppName:  getEnv("APP_NAME", "logistics-delivery"),
		Env:      getEnv("ENV", "development"),
		Version:  getEnv("VERSION", "0.0.1"),
		HTTPPort: getEnv("HTTP_PORT", "8080"),
		GRPCPort: getEnv("GRPC_PORT", "9090"),
		Postgres: PostgresConfig{
			DSN:             getEnv("POSTGRES_DSN", "postgres://logistics:logistics@localhost:5432/logistics?sslmode=disable"),
			MaxConns:        getEnvInt("POSTGRES_MAX_CONNS", 25),
			MinConns:        getEnvInt("POSTGRES_MIN_CONNS", 5),
			MaxConnLifetime: time.Duration(getEnvInt("POSTGRES_MAX_CONN_LIFETIME_MIN", 30)) * time.Minute,
		},
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		Kafka: KafkaConfig{
			Brokers:          splitEnv(getEnv("KAFKA_BROKERS", "localhost:9092")),
			ShipmentTopic:    getEnv("KAFKA_SHIPMENT_TOPIC", "logistics.shipments"),
			TrackingTopic:    getEnv("KAFKA_TRACKING_TOPIC", "logistics.tracking"),
			CourierTopic:     getEnv("KAFKA_COURIER_TOPIC", "logistics.couriers"),
			FulfillmentTopic: getEnv("KAFKA_FULFILLMENT_TOPIC", "logistics.fulfillment"),
			ConsumerGroup:    getEnv("KAFKA_CONSUMER_GROUP", "logistics-consumer-group"),
		},
		OTEL: OTELConfig{
			Endpoint: getEnv("OTEL_EXPORTER_OTLP_ENDPOINT", "localhost:4317"),
			Insecure: getEnv("OTEL_INSECURE", "true") == "true",
		},
	}
}

func splitEnv(s string) []string {
	if s == "" {
		return nil
	}
	return []string{s}
}
