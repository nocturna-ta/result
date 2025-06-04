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
	myRouter.POST("/vote-result", api.InsertVoteResult, router.MustAuthorized(false))

	return myRouter
}
