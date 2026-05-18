package config
import ("os"; "strconv"; "time")
type Config struct { AppName string; AppEnv string; LogLevel string; HTTPPort int; GRPCPort int; Redis RedisConfig; Kafka KafkaConfig; OpenTelemetry OTELConfig }
type RedisConfig struct { Addr string; Password string; DB int; PoolSize int; MinIdleConns int; DialTimeout time.Duration; ReadTimeout time.Duration; WriteTimeout time.Duration; MaxRetries int }
type KafkaConfig struct { Brokers []string }
type OTELConfig struct { Endpoint string; ServiceName string; TraceRatio float64 }
func Load() *Config {
	return &Config{AppName: "shopee-advertising", AppEnv: "development", LogLevel: "info", HTTPPort: 8080, GRPCPort: 9090,
		Redis: RedisConfig{Addr: "localhost:6379", DB: 0, PoolSize: 100, MinIdleConns: 20, DialTimeout: 5 * time.Second, ReadTimeout: 3 * time.Second, WriteTimeout: 3 * time.Second, MaxRetries: 3},
		Kafka: KafkaConfig{Brokers: []string{"localhost:9092"}},
		OpenTelemetry: OTELConfig{Endpoint: "http://localhost:4318", ServiceName: "shopee-advertising", TraceRatio: 0.1}}
}
