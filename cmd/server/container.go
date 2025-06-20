package server

import (
	"context"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/database/sql"
	"github.com/nocturna-ta/golib/ethereum"
	"github.com/nocturna-ta/result/config"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
	"github.com/nocturna-ta/result/internal/interfaces/dao"
	"github.com/nocturna-ta/result/internal/usecases"
	"github.com/nocturna-ta/result/internal/usecases/live_result"
	"time"
)

type container struct {
	Cfg          config.MainConfig
	LiveResultUc usecases.LiveResultUseCases
	WebSocketHub *websocket.Hub
}

type options struct {
	Cfg    *config.MainConfig
	DB     *sql.Store
	Client ethereum.Client
	Ctx    context.Context
	Cache  cache.Cache
}

func newContainer(opts *options) *container {
	voteResultRepo := dao.NewVoteResultRepository(&dao.OptsVoteResultRepository{
		DB: opts.DB,
	})

	liveResultRepo := dao.NewLiveResultRepository(&dao.OptsLiveResultRepository{
		DB: opts.DB,
	})

	wsHub := websocket.NewHub(opts.Ctx)

	liveResultUc := live_result.New(&live_result.Options{
		LiveResultRepo: liveResultRepo,
		VoteResultRepo: voteResultRepo,
		RedisCache:     opts.Cache,
		WsHub:          wsHub,
	})

	fastLiveResultUc := live_result.New(&live_result.Options{
		VoteResultRepo: voteResultRepo,
		LiveResultRepo: liveResultRepo,
		WsHub:          wsHub,
		RedisCache:     opts.Cache,
	})

	go wsHub.Run()

	go liveResultUc.StartIncrementalBroadcast(opts.Ctx, 30*time.Second)

	activeElections := []string{
		"a1234567-bb0b-4103-84f5-edd74e6e1234", "b2345678-bb0b-4103-84f5-edd74e6e2345", "c976902f-bb0b-4103-84f5-edd74e6e928f",
	}

	fastLiveResultUc.StartCacheWarming(opts.Ctx, activeElections)

	fastLiveResultUc.StartIncrementalBroadcast(opts.Ctx, 5*time.Second)

	return &container{
		Cfg:          *opts.Cfg,
		LiveResultUc: liveResultUc,
		WebSocketHub: wsHub,
	}
}
