package repository

import (
	"context"
	"github.com/nocturna-ta/result/internal/domain/model"
)

type VoteResultRepository interface {
	InsertVoteResult(ctx context.Context, result *model.VoteResult) error
	UpdateVoteResult(ctx context.Context, result *model.VoteResult) error
	GetVoteResultByID(ctx context.Context, id string) (*model.VoteResult, error)
}
