package consumer

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	event2 "github.com/nocturna-ta/common-model/models/event"
	libCtx "github.com/nocturna-ta/golib/context"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/domain/model"
	"time"
)

func (m *Module) ConsumeVoteProcessed(ctx context.Context, message *event.EventConsumeMessage) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ConsumerUseCases.ConsumeVoteProcessed")
	defer span.End()

	requestId := libCtx.ReadRequestId(ctx)
	log.WithFields(log.Fields{
		"request_id": requestId,
		"topic":      message.Topic,
		"key":        message.Key,
	}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Processing message")

	var voteMessage event2.VoteProcessedMessage
	err := json.Unmarshal(message.Data, &voteMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"topic":      message.Topic,
			"data":       string(message.Data),
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Failed to unmarshal message")
		return nil
	}

	result := model.VoteResult{
		ID:              uuid.New(),
		VoteID:          uuid.MustParse(voteMessage.VoteID),
		VoterID:         voteMessage.VoterID,
		ElectionPairID:  uuid.New(),
		Region:          "Bandung",
		Status:          voteMessage.Status,
		TransactionHash: voteMessage.TransactionHash,
		ErrorMessage:    voteMessage.ErrorMessage,
		VotedAt:         time.Now(),
		ProcessedAt:     &voteMessage.ProcessedAt,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = m.resultRepo.InsertVoteResult(ctx, &result)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"result":     result,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Failed to insert vote result")
		return err
	}

	return nil
}
