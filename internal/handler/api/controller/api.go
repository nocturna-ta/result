package controller

import (
	"github.com/gofiber/swagger"
	"github.com/nocturna-ta/golib/router"
	_ "github.com/nocturna-ta/result/docs"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
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
	liveResultUc   usecases.LiveResultUseCases
	wsController   *WebSocketController
}

type Options struct {
	Prefix         string
	Port           uint
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	RequestTimeout time.Duration
	EnableSwagger  bool
	LiveResultUc   usecases.LiveResultUseCases
	WebSocketHub   *websocket.Hub
}

func New(opts *Options) *API {

	wsHandler := websocket.NewHandler(opts.WebSocketHub)

	wsController := NewWebSocketController(&WebSocketControllerOptions{
		Handler:      wsHandler,
		LiveResultUc: opts.LiveResultUc,
	})

	return &API{
		prefix:         opts.Prefix,
		port:           opts.Port,
		readTimeout:    opts.ReadTimeout,
		writeTimeout:   opts.WriteTimeout,
		requestTimeout: opts.RequestTimeout,
		enableSwagger:  opts.EnableSwagger,
		liveResultUc:   opts.LiveResultUc,
		wsController:   wsController,
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

		myRouter.CustomHandler("GET", "/result/docs/*", swagger.New(swaggerConfig), router.MustAuthorized(false))
	}

	myRouter.GET("/health", api.Ping, router.MustAuthorized(false))
	myRouter.Group("/v1", func(v1 *router.FastRouter) {
		v1.Group("/live", func(live *router.FastRouter) {
			live.GET("/elections/:election_pair_id", api.GetFastLiveElectionResults, router.MustAuthorized(false))
			live.GET("/elections/:election_pair_id/summary", api.GetFastElectionSummary, router.MustAuthorized(false))
			live.GET("/cities/:city_name", api.GetFastCityResults, router.MustAuthorized(false))
			live.GET("/elections", api.GetFastElectionSummaries, router.MustAuthorized(false))
			live.GET("/rankings/:election_pair_id", api.GetFastCityRankings, router.MustAuthorized(false))

			live.GET("/status", api.wsController.GetLiveResultsStatus, router.MustAuthorized(false))
			live.POST("/broadcast", api.wsController.TriggerBroadcast, router.MustAuthorized(false))

			live.Use("/ws", api.wsController.WebSocketMiddleware())
			live.CustomHandler("GET", "/ws", api.wsController.HandleWebSocket, router.MustAuthorized(false))
		})
	})

	return myRouter
}
