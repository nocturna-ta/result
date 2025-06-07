package consumer

import (
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/event/handler"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/interfaces/dao"
	"github.com/nocturna-ta/result/internal/usecases"
	"github.com/nocturna-ta/result/internal/usecases/consumer"
)

type container struct {
	Cfg          config.MainConfig
	ConsumerUc   usecases.Consumer
	EventHandler handler.EventHandler
}

type options struct {
	Cfg *config.MainConfig
	DB  *sql.Store
}

func newContainer(opts *options) *container {
	resultRepo := dao.NewVoteResultRepository(&dao.OptsVoteResultRepository{
		DB: opts.DB,
	})

	consumerUc := consumer.New(&consumer.Options{
		ResultRepo: resultRepo,
		Topics:     opts.Cfg.Kafka.Topics,
	})

	eventHandler := handler.New(&handler.Options{
		RetryConfig: handler.RetryConfig{
			MaxRetry:          opts.Cfg.Kafka.Consumer.Retry.MaxRetry,
			RetryInitialDelay: opts.Cfg.Kafka.Consumer.Retry.RetryInitialDelay,
			MaxJitter:         opts.Cfg.Kafka.Consumer.Retry.MaxJitter,
			HandlerTimeout:    opts.Cfg.Kafka.Consumer.Retry.HandlerTimeout,
			BackOffConfig:     opts.Cfg.Kafka.Consumer.Retry.BackOffConfig,
		},
		Publisher:   nil,
		DlqTopic:    opts.Cfg.Kafka.Topics.VoteDLQ.Value,
		ServiceName: "result-service",
	})

	return &container{
		Cfg:          *opts.Cfg,
		ConsumerUc:   consumerUc,
		EventHandler: eventHandler,
	}
}
