package fast_live_result

import (
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/result/internal/domain/repository"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
	"github.com/nocturna-ta/result/internal/usecases"
)

type CacheStatistics struct {
	TotalHits   uint64
	TotalMisses uint64
}

type Module struct {
	liveResultRepo  repository.LiveResultRepository
	voteResultRepo  repository.VoteResultRepository
	redisCache      cache.Cache
	wsHub           *websocket.Hub
	cacheStatistics *CacheStatistics
}

type Options struct {
	LiveResultRepo repository.LiveResultRepository
	VoteResultRepo repository.VoteResultRepository
	RedisCache     cache.Cache
	WsHub          *websocket.Hub
}

func New(opts *Options) usecases.FastLiveResultUseCases {
	return &Module{
		liveResultRepo:  opts.LiveResultRepo,
		voteResultRepo:  opts.VoteResultRepo,
		redisCache:      opts.RedisCache,
		wsHub:           opts.WsHub,
		cacheStatistics: &CacheStatistics{},
	}
}
