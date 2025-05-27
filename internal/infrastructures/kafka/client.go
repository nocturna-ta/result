package kafka

import (
	"context"
	"github.com/nocturna-ta/golib/event"
	_ "github.com/nocturna-ta/golib/event/kafka"
	"github.com/nocturna-ta/golib/tracing/newrelic"

	"github.com/nocturna-ta/result/config"
)

func NewPublisher(ctx context.Context, config config.KafkaProducerConfig) (event.MessagePublisher, error) {
	conf := &event.PublisherConfig{
		DriverConfig: &event.DriverConfig{
			Type: "kafka",
			Config: map[string]any{
				"brokers":      config.Brokers,
				"max_attempts": config.MaxAttempt,
				"idempotent":   config.Idempotent,
			},
		},
	}

	return event.NewPublisher(ctx, conf)
}

func NewConsumer(ctx context.Context, config config.KafkaConsumerConfig, eventHandler event.ConsumerEventHandler) (*event.Consumer, error) {
	conf := &event.ConsumerConfig{
		Consumer: &event.DriverConfig{
			Type: "kafka",
			Config: map[string]any{
				"brokers":               config.Brokers,
				"kafka_cluster_version": config.ClusterVersion,
			},
		},
		CommitStrategy: &event.DriverConfig{
			Type:   "commit_on_success",
			Config: map[string]any{},
		},
		EventHandler: nil,
		NewRelicOpts: &newrelic.Options{
			MustStart: false,
		},
	}

	if eventHandler != nil {
		conf.EventHandler = eventHandler
	}

	return event.NewConsumer(ctx, conf)
}
