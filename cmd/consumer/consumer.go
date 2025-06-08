package consumer

import (
	"context"
	"github.com/nocturna-ta/golib/database/sql"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	configLocation, _ := cmd.Flags().GetString("config")
	cfg := &config.MainConfig{}
	config.ReadConfig(cfg, configLocation)

	database := sql.New(sql.DBConfig{
		SlaveDSN:        cfg.Database.SlaveDSN,
		MasterDSN:       cfg.Database.MasterDSN,
		RetryInterval:   cfg.Database.RetryInterval,
		MaxIdleConn:     cfg.Database.MaxIdleConn,
		MaxConn:         cfg.Database.MaxConn,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}, sql.DriverClickHouse)

	appContainer := newContainer(&options{
		Cfg: cfg,
		DB:  database,
		Ctx: ctx,
	})

	consumer, err := kafka.NewConsumer(context.Background(), cfg.Kafka.Consumer, &appContainer.EventHandler)
	if err != nil {
		log.Fatalf("Failed to instantiate kafka consumer: %w", err)
	}

	topicHandler := map[event.TopicName]event.ConsumerHandlerConfig{
		event.TopicName(cfg.Kafka.Topics.VoteSubmitData.Value): {
			ConsumerGroup:     cfg.Kafka.Consumer.ConsumerGroup,
			ErrorHandlerLevel: cfg.Kafka.Topics.VoteSubmitData.ErrorHandler,
			Handler:           appContainer.ConsumerUc.ConsumeVoteSubmit,
			WithBackOff:       cfg.Kafka.Topics.VoteSubmitData.WithBackOff,
		},
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
