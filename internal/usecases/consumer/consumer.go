package consumer

import (
	"context"
	"encoding/json"
	event2 "github.com/nocturna-ta/common-model/models/event"
	libCtx "github.com/nocturna-ta/golib/context"
	"github.com/nocturna-ta/golib/event"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/pkg/constants"
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

	existingResult, err := m.resultRepo.GetVoteResultByID(ctx, voteMessage.VoteID)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"vote_id":    voteMessage.VoteID,
		}).DebugWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Failed to get existing vote result")
	}

	var electionPairID, region string

	if existingResult != nil {
		existingResult.Status = voteMessage.Status
		existingResult.TransactionHash = voteMessage.TransactionHash
		existingResult.ErrorMessage = voteMessage.ErrorMessage
		existingResult.ProcessedAt = &voteMessage.ProcessedAt

		electionPairID = existingResult.ElectionPairID
		region = existingResult.Region

		err = m.resultRepo.UpdateVoteResult(ctx, existingResult)
		if err != nil {
			log.WithFields(log.Fields{
				"request_id": requestId,
				"error":      err,
				"result":     existingResult,
			}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Failed to update existing vote result")
			return err
		}
		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteMessage.VoteID,
			"status":     voteMessage.Status,
		}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Updated existing vote result")
	} else {
		result := model.FromVoteProcessedMessage(&voteMessage)
		electionPairID = result.ElectionPairID
		region = result.Region

		err = m.resultRepo.InsertVoteResult(ctx, result)
		if err != nil {
			log.WithFields(log.Fields{
				"request_id": requestId,
				"error":      err,
				"result":     result,
			}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Failed to insert new vote result")
			return err
		}

		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteMessage.VoteID,
			"status":     voteMessage.Status,
		}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumeVoteProcessed] Inserted new vote result")
	}

	m.broadcastLiveUpdates(ctx, voteMessage.VoteID, electionPairID, region)

	return nil
}

func (m *Module) ConsumeVoteSubmit(ctx context.Context, message *event.EventConsumeMessage) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "ConsumerUseCases.ConsumerVoteSubmit")
	defer span.End()

	requestId := libCtx.ReadRequestId(ctx)
	log.WithFields(log.Fields{
		"request_id": requestId,
		"topic":      message.Topic,
		"key":        message.Key,
	}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Processing message")

	var voteMessage event2.VoteSubmitMessage
	err := json.Unmarshal(message.Data, &voteMessage)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"topic":      message.Topic,
			"data":       string(message.Data),
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Failed to unmarshal message")
		return nil
	}

	operation, ok := message.Metadata[constants.MetaDataOperation].(string)
	if !ok {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      "missing operation metadata",
			"topic":      message.Topic,
			"metadata":   message.Metadata,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Missing operation metadata")
		operation = constants.Create
	}

	var electionPairID, region string

	switch operation {
	case constants.Create:
		electionPairID, region = m.handleVoteCreate(ctx, &voteMessage, requestId)
	case constants.Update:
		electionPairID, region = m.handleVoteUpdate(ctx, &voteMessage, requestId)
	default:
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      "unknown operation",
			"operation":  operation,
			"topic":      message.Topic,
		}).WarnWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Unknown operation, skipping")
		return nil
	}

	m.broadcastLiveUpdates(ctx, voteMessage.VoteID, electionPairID, region)

	return nil
}

func (m *Module) handleVoteCreate(ctx context.Context, voteMessage *event2.VoteSubmitMessage, requestId string) (string, string) {
	existingResult, err := m.resultRepo.GetVoteResultByID(ctx, voteMessage.VoteID)
	if err == nil && existingResult != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteMessage.VoteID,
		}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Vote already exists, skipping creation")
		return existingResult.ElectionPairID, existingResult.Region
	}

	result := model.FromVoteSubmitMessage(voteMessage)
	err = m.resultRepo.InsertVoteResult(ctx, result)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"result":     result,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Failed to insert new vote result")
		return "", ""
	}

	log.WithFields(log.Fields{
		"request_id":       requestId,
		"vote_id":          voteMessage.VoteID,
		"election_pair_id": voteMessage.ElectionPairID,
		"region":           voteMessage.Region,
	}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Inserted new vote result")

	return voteMessage.ElectionPairID, voteMessage.Region
}

func (m *Module) handleVoteUpdate(ctx context.Context, voteMessage *event2.VoteSubmitMessage, requestId string) (string, string) {
	existingResult, err := m.resultRepo.GetVoteResultByID(ctx, voteMessage.VoteID)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"vote_id":    voteMessage.VoteID,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Failed to get existing vote result")
		return "", ""
	}
	if existingResult == nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"vote_id":    voteMessage.VoteID,
		}).WarnWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Vote does not exist, skipping update")
		return "", ""
	}

	existingResult.ElectionPairID = voteMessage.ElectionPairID
	existingResult.Region = voteMessage.Region
	existingResult.VotedAt = voteMessage.SubmittedAt

	err = m.resultRepo.UpdateVoteResult(ctx, existingResult)
	if err != nil {
		log.WithFields(log.Fields{
			"request_id": requestId,
			"error":      err,
			"result":     existingResult,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Failed to update existing vote result")
		return "", ""
	}
	log.WithFields(log.Fields{
		"request_id": requestId,
		"vote_id":    voteMessage.VoteID,
	}).InfoWithCtx(ctx, "[ConsumerUseCases.ConsumerVoteSubmit] Updated existing vote result")

	return voteMessage.ElectionPairID, voteMessage.Region
}

func (m *Module) broadcastLiveUpdates(ctx context.Context, voteID, electionPairID, region string) {
	if m.liveResult == nil || m.liveResult.GetConnectedClients(ctx) == 0 {
		return
	}

	if err := m.liveResult.BroadcastVoteUpdate(ctx, voteID); err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"vote_id": voteID,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases] Failed to broadcast vote update")
	}

	if err := m.liveResult.BroadcastAllUpdates(ctx, electionPairID, region); err != nil {
		log.WithFields(log.Fields{
			"error":            err,
			"election_pair_id": electionPairID,
			"region":           region,
		}).ErrorWithCtx(ctx, "[ConsumerUseCases] Failed to broadcast aggregated updates")
	}
}
