package config

import (
	"github.com/nocturna-ta/golib/config"
	"github.com/nocturna-ta/golib/log"
	"time"
)

type (
	MainConfig struct {
		Server     ServerConfig     `yaml:"Server"`
		API        APIConfig        `yaml:"API"`
		Database   DBConfig         `yaml:"Database"`
		Blockchain BlockchainConfig `yaml:"BlockchainConfig"`
		JWT        JWTConfig        `yaml:"JWT"`
		Kafka      KafkaConfig      `yaml:"Kafka"`
		Encryption EncryptionConfig `yaml:"Encryption"`
		Cors       CorsConfig       `yaml:"Cors"`
		GrpcServer GrpcServerConfig `yaml:"GrpcServer"`
	}

	ServerConfig struct {
		Port         uint          `yaml:"Port" env:"SERVER_PORT"`
		WriteTimeout time.Duration `yaml:"WriteTimeout" env:"SERVER_WRITE_TIMEOUT"`
		ReadTimeout  time.Duration `yaml:"ReadTimeout" env:"SERVER_READ_TIMEOUT"`
	}

	APIConfig struct {
		BasePath      string        `yaml:"BasePath" env:"API_BASE_PATH"`
		APITimeout    time.Duration `yaml:"APITimeout" env:"API_TIMEOUT"`
		EnableSwagger bool          `yaml:"EnableSwagger" env:"ENABLE_SWAGGER" default:"false"`
	}

	DBConfig struct {
		SlaveDSN        string `yaml:"SlaveDSN" env:"DB_SLAVE_DSN"`
		MasterDSN       string `yaml:"MasterDSN" env:"DB_MASTER_DSN"`
		RetryInterval   int    `yaml:"RetryInterval" env:"DB_RETRY_INTERVAL"`
		MaxIdleConn     int    `yaml:"MaxIdleConn" env:"DB_MAX_IDLE_CONN"`
		MaxConn         int    `yaml:"MaxConn" env:"DB_MAX_CONN"`
		ConnMaxLifetime string `yaml:"ConnMaxLifetime" env:"DB_CONN_MAX_LIFETIME"`
	}

	BlockchainConfig struct {
		GanacheURL             string `yaml:"GanacheURL"`
		VotechainAddress       string `yaml:"VotechainAddress" `
		VotechainBaseAddress   string `yaml:"VotechainBaseAddress"`
		KPUManagerAddress      string `yaml:"KPUManagerAddress"`
		VoterManagerAddress    string `yaml:"VoterManagerAddress"`
		ElectionManagerAddress string `yaml:"ElectionManagerAddress"`
	}

	JWTConfig struct {
		Secret string `yaml:"Secret" env:"JWT_SECRET"`
	}

	EncryptionConfig struct {
		Key string `yaml:"Key" env:"ENCRYPTION_KEY"`
	}

	CorsConfig struct {
		AllowOrigins     string `yaml:"AllowOrigins"`
		AllowMethods     string `yaml:"AllowMethods"`
		AllowHeaders     string `yaml:"AllowHeaders"`
		AllowCredentials bool   `yaml:"AllowCredentials"`
		ExposeHeaders    string `yaml:"ExposeHeaders"`
		MaxAge           int    `yaml:"MaxAge"`
	}

	GrpcServerConfig struct {
		Port uint `yaml:"Port"`
	}

	KafkaConfig struct {
		Producer KafkaProducerConfig `yaml:"Producer"`
		Consumer KafkaConsumerConfig `yaml:"Consumer"`
		Topics   KafkaTopics         `yaml:"Topics"`
	}

	KafkaProducerConfig struct {
		Brokers    []string `yaml:"Brokers" env:"KAFKA_BROKERS"`
		Idempotent bool     `yaml:"Idempotent" env:"KAFKA_IDEMPOTENT"`
		MaxAttempt int      `yaml:"MaxAttempt" env:"KAFKA_MAX_ATTEMPTS"`
	}

	KafkaConsumerConfig struct {
		Brokers        []string    `yaml:"Brokers"`
		ClusterVersion string      `yaml:"ClusterVersion"`
		ConsumerGroup  string      `yaml:"ConsumerGroup"`
		MaxRetries     int         `yaml:"MaxRetries"`
		WorkerPoolSize int         `yaml:"WorkerPoolSize"`
		MaxAttempt     int         `yaml:"MaxAttempt"`
		Retry          RetryConfig `yaml:"Retry"`
	}

	RetryConfig struct {
		MaxRetry          int             `yaml:"MaxRetry"`
		RetryInitialDelay time.Duration   `yaml:"RetryInitialDelay"`
		MaxJitter         time.Duration   `yaml:"MaxJitter"`
		HandlerTimeout    time.Duration   `yaml:"HandlerTimeout"`
		BackOffConfig     []time.Duration `yaml:"BackOffConfig"`
	}

	KafkaTopics struct {
		VoteSubmitData KafkaTopicConfig `yaml:"VoteSubmitData"`
		VoteProcessed  KafkaTopicConfig `yaml:"VoteProcessed"`
		VoteDLQ        KafkaTopicConfig `yaml:"VoteDLQ"`
	}

	KafkaTopicConfig struct {
		Value        string `yaml:"Value" env:"KAFKA_TOPIC_VALUE"`
		ErrorHandler string `yaml:"ErrorHandler"`
		WithBackOff  bool   `yaml:"WithBackOff"`
	}
)

func ReadConfig(cfg any, configLocation string) {
	if configLocation == "" {
		configLocation = "file://config/files/config.yaml"
	}

	if err := config.ReadConfig(cfg, configLocation, true); err != nil {
		log.WithFields(log.Fields{
			"error":           err,
			"config-location": configLocation,
		}).Fatal("Failed to read config")
	}
}
