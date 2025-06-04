package consumer

import (
	"context"
	"github.com/nocturna-ta/golib/database/nosql/clickhouse"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/infrastructures/kafka"
	"github.com/spf13/cobra"
)

var (
	serveConsumerCmd = &cobra.Command{
		Use:   "run-consumer",
		Short: "Result Consumer Service",
		RunE:  run,
	}
)

func ServeConsumerCmd() *cobra.Command {
	serveConsumerCmd.Flags().StringP("config", "c", "", "Config Path, both relative or absolute. i.e: /usr/local/bin/config/files")
	return serveConsumerCmd
}

func run(cmd *cobra.Command, args []string) error {
	configLocation, _ := cmd.Flags().GetString("config")

	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database, err := clickhouse.New(&clickhouse.Config{
		Addrs: cfg.ClickHouse.Addrs,
		Auth: clickhouse.Auth{
			Database: cfg.ClickHouse.Auth.Database,
			Username: cfg.ClickHouse.Auth.Username,
			Password: cfg.ClickHouse.Auth.Password,
		},
		Database:        cfg.ClickHouse.Database,
		DialTimeout:     cfg.ClickHouse.DialTimeout,
		MaxOpenConns:    cfg.ClickHouse.MaxOpenConns,
		MaxIdleConns:    cfg.ClickHouse.MaxIdleConns,
		ConnMaxLifetime: cfg.ClickHouse.ConnMaxLifetime,
		TLS: &clickhouse.TLSConfig{
			Enable:             cfg.ClickHouse.TLS.Enable,
			InsecureSkipVerify: cfg.ClickHouse.TLS.InsecureSkipVerify,
			CertFile:           cfg.ClickHouse.TLS.CertFile,
			KeyFile:            cfg.ClickHouse.TLS.KeyFile,
			CAFile:             cfg.ClickHouse.TLS.CAFile,
		},
		BlockBufferSize:    cfg.ClickHouse.BlockBufferSize,
		MaxCompressionSize: cfg.ClickHouse.MaxCompressionSize,
		AsyncInsert:        cfg.ClickHouse.AsyncInsert,
		AsyncInsertOptions: clickhouse.AsyncInsertOptions{
			MaxBatchSize: cfg.ClickHouse.AsyncInsertOptions.MaxBatchSize,
			MaxDelay:     cfg.ClickHouse.AsyncInsertOptions.MaxDelay,
		},
		Debug: cfg.ClickHouse.Debug,
	})

	log.Info("ClickHouse connections established successfully")

	appContainer := newContainer(&options{
		Cfg: cfg,
		DB:  database,
	})

	consumer, err := kafka.NewConsumer(context.Background(), cfg.Kafka.Consumer, &appContainer.EventHandler)
	if err != nil {
		log.Fatalf("Failed to instantiate kafka consumer: %w", err)
	}

	topicHandler := map[event.TopicName]event.ConsumerHandlerConfig{
		event.TopicName(cfg.Kafka.Topics.VoteProcessed.Value): {
			ConsumerGroup:     cfg.Kafka.Consumer.ConsumerGroup,
			ErrorHandlerLevel: cfg.Kafka.Topics.VoteProcessed.ErrorHandler,
			Handler:           appContainer.ConsumerUc.ConsumeVoteProcessed,
			WithBackOff:       cfg.Kafka.Topics.VoteProcessed.WithBackOff,
		},
	}

	consumer.RunWithHandlerConfig(topicHandler)

	return nil
}
