package consumer

import (
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/domain/repository"
	"github.com/nocturna-ta/result/internal/usecases"
)

type Module struct {
	resultRepo repository.VoteResultRepository
	topics     config.KafkaTopics
}

type Options struct {
	ResultRepo repository.VoteResultRepository
	Topics     config.KafkaTopics
}

func New(opts *Options) usecases.Consumer {
	return &Module{
		resultRepo: opts.ResultRepo,
		topics:     opts.Topics,
	}
}
