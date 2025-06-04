package usecases

import (
	"context"
	"github.com/nocturna-ta/result/internal/usecases/request"
	"github.com/nocturna-ta/result/internal/usecases/response"
)

type VoteResultUseCases interface {
	InsertVoteResult(ctx context.Context, entry request.VoteResultEntry) (*response.EntryResponse, error)
}
