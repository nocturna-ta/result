package server

import (
	"context"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
	"github.com/nocturna-ta/result/internal/interfaces/dao"
	"github.com/nocturna-ta/result/internal/usecases"
	"github.com/nocturna-ta/result/internal/usecases/live_result"
	"github.com/nocturna-ta/result/internal/usecases/vote_result"
	"time"
)

type container struct {
	Cfg          config.MainConfig
	VoteResultUc usecases.VoteResultUseCases
	LiveResultUc usecases.LiveResultUsecases
	WebSocketHub *websocket.Hub
}

type options struct {
	Cfg    *config.MainConfig
	DB     *sql.Store
	Client ethereum.Client
	Ctx    context.Context
}

func newContainer(opts *options) *container {
	voteResultRepo := dao.NewVoteResultRepository(&dao.OptsVoteResultRepository{
		DB: opts.DB,
	})

	voteResultUc := vote_result.New(&vote_result.Opts{
		VoteResultRepo: voteResultRepo,
	})

	wsHub := websocket.NewHub(opts.Ctx)

	liveResultUc := live_result.New(&live_result.Options{
		VoteResultRepo: voteResultRepo,
		Hub:            wsHub,
	})

	go wsHub.Run()

	go liveResultUc.StartPeriodicBroadcast(opts.Ctx, 30*time.Second)

	return &container{
		Cfg:          *opts.Cfg,
		VoteResultUc: voteResultUc,
		LiveResultUc: liveResultUc,
		WebSocketHub: wsHub,
	}
}
