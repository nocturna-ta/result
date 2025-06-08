package live_result

import (
	"context"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

func (m *Module) BroadcastVoteUpdate(ctx context.Context, voteID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultUseCase.BroadcastVoteUpdate")
	defer span.End()

	voteResult, err := m.voteResultRepo.GetVoteResultByID(ctx, voteID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"voteID": voteID,
		}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastVoteUpdate] Failed to get vote result by ID")
		return err
	}

	voteResponse := &response.VoteResultResponse{
		ID:              voteResult.ID,
		VoterID:         voteResult.VoterID,
		ElectionPairID:  voteResult.ElectionPairID,
		Region:          voteResult.Region,
		Status:          voteResult.Status,
		TransactionHash: voteResult.TransactionHash,
		ErrorMessage:    voteResult.ErrorMessage,
		VotedAt:         voteResult.VotedAt,
		ProcessedAt:     voteResult.ProcessedAt,
		CreatedAt:       voteResult.CreatedAt,
		UpdatedAt:       voteResult.UpdatedAt,
	}

	m.hub.BroadcastVoteUpdate(voteResponse)

	log.WithFields(log.Fields{
		"vote_id":          voteID,
		"election_pair_id": voteResult.ElectionPairID,
		"region":           voteResult.Region,
		"status":           voteResult.Status,
		"clients":          m.hub.GetClientCount(),
	}).InfoWithCtx(ctx, "[LiveResultService.BroadcastVoteUpdate] Vote update broadcasted")

	return nil
}

func (m *Module) BroadcastElectionUpdate(ctx context.Context, electionPairID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultUseCase.BroadcastElectionUpdate")
	defer span.End()

	electionResult, err := m.voteResultRepo.GetElectionResults(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"electionPairID": electionPairID,
		}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastElectionUpdate] Failed to get election results")
		return err
	}

	electionResponse := &response.ElectionVoteResultResponse{
		ElectionPairID: electionResult.ElectionPairID,
		Region:         electionResult.Region,
		TotalVotes:     electionResult.TotalVotes,
		ConfirmedVotes: electionResult.ConfirmedVotes,
		PendingVotes:   electionResult.PendingVotes,
		ErrorVotes:     electionResult.ErrorVotes,
		LastUpdated:    electionResult.LastUpdated,
	}

	m.hub.BroadcastElectionUpdate(electionResponse)

	log.WithFields(log.Fields{
		"election_pair_id": electionPairID,
		"total_votes":      electionResult.TotalVotes,
		"confirmed_votes":  electionResult.ConfirmedVotes,
		"clients":          m.hub.GetClientCount(),
	}).InfoWithCtx(ctx, "[LiveResultService.BroadcastElectionUpdate] Election update broadcasted")

	return nil
}

func (m *Module) BroadcastRegionUpdate(ctx context.Context, region string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultUseCase.BroadcastRegionUpdate")
	defer span.End()

	regionResult, err := m.voteResultRepo.GetRegionResults(ctx, region)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastRegionUpdate] Failed to get region results")
		return err
	}

	regionResponse := &response.RegionVoteResultResponse{
		Region:         regionResult.Region,
		TotalVotes:     regionResult.TotalVotes,
		ConfirmedVotes: regionResult.ConfirmedVotes,
		PendingVotes:   regionResult.PendingVotes,
		ErrorVotes:     regionResult.ErrorVotes,
		LastUpdated:    regionResult.LastUpdated,
	}

	m.hub.BroadcastRegionUpdate(regionResponse)

	log.WithFields(log.Fields{
		"region":          region,
		"total_votes":     regionResult.TotalVotes,
		"confirmed_votes": regionResult.ConfirmedVotes,
		"clients":         m.hub.GetClientCount(),
	}).InfoWithCtx(ctx, "[LiveResultService.BroadcastRegionUpdate] Region update broadcasted")

	return nil

}

func (m *Module) BroadcastStatisticsUpdate(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultUseCase.BroadcastStatisticsUpdate")
	defer span.End()

	stats, err := m.voteResultRepo.GetOverallStatistics(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastStatisticsUpdate] Failed to get overall statistics")
		return err
	}

	stats.CalculateSuccessRate()

	statsResponse := &response.VoteStatisticsResponse{
		TotalVotes:     stats.TotalVotes,
		ConfirmedVotes: stats.ConfirmedVotes,
		PendingVotes:   stats.PendingVotes,
		ErrorVotes:     stats.ErrorVotes,
		SuccessRate:    stats.SuccessRate,
		LastUpdated:    stats.LastUpdated,
	}
	m.hub.BroadcastStatisticsUpdate(statsResponse)

	log.WithFields(log.Fields{
		"total_votes":     stats.TotalVotes,
		"confirmed_votes": stats.ConfirmedVotes,
		"success_rate":    stats.SuccessRate,
		"clients":         m.hub.GetClientCount(),
	}).InfoWithCtx(ctx, "[LiveResultService.BroadcastStatisticsUpdate] Statistics update broadcasted")

	return nil
}

func (m *Module) BroadcastAllUpdates(ctx context.Context, electionPairID, region string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "LiveResultUseCase.BroadcastAllUpdates")
	defer span.End()

	if electionPairID != "" {
		if err := m.BroadcastElectionUpdate(ctx, electionPairID); err != nil {
			log.WithFields(log.Fields{
				"error":          err,
				"electionPairID": electionPairID,
			}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastAllUpdates] Failed to broadcast election update")
		}
	}
	if region != "" {
		if err := m.BroadcastRegionUpdate(ctx, region); err != nil {
			log.WithFields(log.Fields{
				"error":  err,
				"region": region,
			}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastAllUpdates] Failed to broadcast region update")
		}
	}
	if err := m.BroadcastStatisticsUpdate(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[LiveResultUseCase.BroadcastAllUpdates] Failed to broadcast statistics update")
	}

	return nil
}

func (m *Module) GetConnectedClients(ctx context.Context) int {
	return m.hub.GetClientCount()
}

func (m *Module) StartPeriodicBroadcast(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	log.WithFields(log.Fields{
		"interval": interval,
	}).InfoWithCtx(ctx, "[LiveResultUseCase.StartPeriodicBroadcast] Starting periodic broadcast")

	for {
		select {
		case <-ctx.Done():
			log.InfoWithCtx(ctx, "[LiveResultUseCase.StartPeriodicBroadcast] Stopping periodic broadcast")
			return
		case <-ticker.C:
			if m.hub.GetClientCount() > 0 {
				if err := m.BroadcastStatisticsUpdate(ctx); err != nil {
					log.WithFields(log.Fields{
						"error": err,
					}).ErrorWithCtx(ctx, "[LiveResultUseCase.StartPeriodicBroadcast] Failed to broadcast statistics update")
				}
			}
		}
	}
}
