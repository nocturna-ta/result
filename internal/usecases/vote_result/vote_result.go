package vote_result

import (
	"context"
	"github.com/google/uuid"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/internal/usecases/request"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

func (m *Module) InsertVoteResult(ctx context.Context, entry request.VoteResultEntry) (*response.EntryResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "VoteResultUseCases.InsertVoteResult")
	defer span.End()

	now := time.Now()

	voteResult := model.VoteResult{
		ID:              uuid.New(),
		VoteID:          uuid.MustParse(entry.VoteID),
		VoterID:         entry.VoterID,
		ElectionPairID:  uuid.MustParse(entry.ElectionPairID),
		Region:          entry.Region,
		Status:          "pending",
		TransactionHash: entry.TransactionHash,
		ErrorMessage:    "n",
		VotedAt:         now,
		ProcessedAt:     &now,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err := m.voteResultRepo.InsertVoteResult(ctx, &voteResult)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"result": voteResult,
		}).ErrorWithCtx(ctx, "[VoteResultUseCases.InsertVoteResult] failed to insert vote result")
		return nil, err
	}

	return &response.EntryResponse{
		ID: voteResult.ID.String(),
	}, nil
}
