package server

import (
	"github.com/nocturna-ta/golib/database/nosql/clickhouse"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/interfaces/dao"
	"github.com/nocturna-ta/result/internal/usecases"
	"github.com/nocturna-ta/result/internal/usecases/vote_result"
)

type container struct {
	Cfg          config.MainConfig
	VoteResultUc usecases.VoteResultUseCases
}

type options struct {
	Cfg       *config.MainConfig
	DB        clickhouse.Client
	Client    ethereum.Client
	Publisher event.MessagePublisher
}

func newContainer(opts *options) *container {
	voteResultRepo := dao.NewVoteResultRepository(&dao.OptsVoteResultRepository{
		DB: opts.DB,
	})

	voteResultUc := vote_result.New(&vote_result.Opts{
		VoteResultRepo: voteResultRepo,
	})

	return &container{
		Cfg:          *opts.Cfg,
		VoteResultUc: voteResultUc,
	}
}
