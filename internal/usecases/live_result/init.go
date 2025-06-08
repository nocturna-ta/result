package live_result

import (
	"github.com/nocturna-ta/result/internal/domain/repository"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
	"github.com/nocturna-ta/result/internal/usecases"
)

type Module struct {
	voteResultRepo repository.VoteResultRepository
	hub            *websocket.Hub
}

type Options struct {
	VoteResultRepo repository.VoteResultRepository
	Hub            *websocket.Hub
}

func New(opts *Options) usecases.LiveResultUsecases {
	return &Module{
		voteResultRepo: opts.VoteResultRepo,
		hub:            opts.Hub,
	}
}
