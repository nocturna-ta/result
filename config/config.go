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
		ClickHouse ClickHouseConfig `yaml:"ClickHouse"`
		Kafka      KafkaConfig      `yaml:"Kafka"`
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

	ClickHouseConfig struct {
		Addrs              []string           `yaml:"Addrs"`
		Auth               Auth               `yaml:"Auth" `
		Database           string             `yaml:"Database" `
		DialTimeout        time.Duration      `yaml:"DialTimeout" `
		MaxOpenConns       int                `yaml:"MaxOpenConns" `
		MaxIdleConns       int                `yaml:"MaxIdleConns" `
		ConnMaxLifetime    time.Duration      `yaml:"ConnMaxLifetime" `
		TLS                *TLSConfig         `yaml:"TLS"`
		BlockBufferSize    uint8              `yaml:"BlockBufferSize"`
		MaxCompressionSize uint64             `yaml:"MaxCompressionSize"`
		AsyncInsert        bool               `yaml:"AsyncInsert"`
		AsyncInsertOptions AsyncInsertOptions `yaml:"AsyncInsertOptions"`
		Debug              bool               `yaml:"Debug"`
	}

	Auth struct {
		Database string `yaml:"Database"`
		Username string `yaml:"Username"`
		Password string `yaml:"Password"`
	}

	TLSConfig struct {
		Enable             bool   `yaml:"Enable"`
		InsecureSkipVerify bool   `yaml:"InsecureSkipVerify"`
		CertFile           string `yaml:"CertFile"`
		KeyFile            string `yaml:"KeyFile"`
		CAFile             string `yaml:"CAFile"`
	}

	AsyncInsertOptions struct {
		MaxBatchSize int           `yaml:"MaxBatchSize"`
		MaxDelay     time.Duration `yaml:"MaxDelay"`
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
		Consumer KafkaConsumerConfig `yaml:"Consumer"`
		Topics   KafkaTopics         `yaml:"Topics"`
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
