package controller

import (
	"github.com/gofiber/swagger"
	"github.com/nocturna-ta/golib/router"
	_ "github.com/nocturna-ta/result/docs"
	"github.com/nocturna-ta/result/internal/usecases"
	"github.com/nocturna-ta/result/pkg/utils"
	"html/template"
	"time"
)

type API struct {
	prefix         string
	port           uint
	readTimeout    time.Duration
	writeTimeout   time.Duration
	requestTimeout time.Duration
	enableSwagger  bool
	voteResult     usecases.VoteResultUseCases
}

type Options struct {
	Prefix         string
	Port           uint
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	EnableSwagger  bool
	VoteResult     usecases.VoteResultUseCases
}

func New(opts *Options) *API {
	return &API{
		prefix:         opts.Prefix,
		port:           opts.Port,
		readTimeout:    opts.ReadTimeout,
		writeTimeout:   opts.WriteTimeout,
		requestTimeout: opts.RequestTimeout,
		enableSwagger:  opts.EnableSwagger,
		voteResult:     opts.VoteResult,
	}
}

func (api *API) RegisterRoute() *router.FastRouter {
	myRouter := router.New(&router.Options{
		Prefix:         api.prefix,
		Port:           api.port,
		ReadTimeout:    api.readTimeout,
		WriteTimeout:   api.writeTimeout,
		RequestTimeout: api.requestTimeout,
	})

	if api.enableSwagger {
		swaggerConfig := swagger.Config{
			Title:        "API Documentation",
			DeepLinking:  true,
			DocExpansion: "list",
			CustomStyle:  template.CSS(utils.ClaudeDarkTheme),
		}

		myRouter.CustomHandler("GET", "/docs/*", swagger.New(swaggerConfig), router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))
	myRouter.Group("/v1", func(v1 *router.FastRouter) {
		v1.Group("/results", func(results *router.FastRouter) {
			results.GET("/votes/:id", api.GetVoteResult, router.MustAuthorized(false))
			results.GET("/votes", api.GetVoteResultByStatus, router.MustAuthorized(false))
			results.GET("/votes/count", api.CountVotesByStatus, router.MustAuthorized(false))

			results.GET("/elections/:election_pair_id", api.GetElectionResults, router.MustAuthorized(false))
			results.GET("/elections/:election_pair_id/votes", api.GetVoteResultByElectionPair, router.MustAuthorized(false))
			results.GET("/elections/:election_pair_id/count", api.CountVotesByElectionPair, router.MustAuthorized(false))

			results.GET("/regions", api.GetRegionStatistics, router.MustAuthorized(false))
			results.GET("/regions/:region", api.GetRegionResults, router.MustAuthorized(false))
			results.GET("/regions/:region/votes", api.GetVoteResultByRegion, router.MustAuthorized(false))
			results.GET("/regions/:region/elections", api.GetElectionResultsByRegion, router.MustAuthorized(false))
			results.GET("/regions/:region/count", api.CountVotesByRegion, router.MustAuthorized(false))

			results.GET("/statistics", api.GetOverallStatistics, router.MustAuthorized(false))
			results.GET("/statistics/daily", api.GetDailyStatistics, router.MustAuthorized(false))
		})
	})

	return myRouter
}
