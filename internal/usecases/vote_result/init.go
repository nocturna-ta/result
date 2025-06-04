package vote_result

import (
	"github.com/nocturna-ta/result/internal/domain/repository"
	"github.com/nocturna-ta/result/internal/usecases"
)

type Module struct {
	voteResultRepo repository.VoteResultRepository
}

type Opts struct {
	VoteResultRepo repository.VoteResultRepository
}

func New(opts *Opts) usecases.VoteResultUseCases {
	return &Module{
		voteResultRepo: opts.VoteResultRepo,
	}
}
