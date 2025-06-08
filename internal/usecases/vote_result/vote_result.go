package vote_result

import (
	"context"
	"errors"
	"github.com/nocturna-ta/golib/custerr"
	"github.com/nocturna-ta/golib/log"
	response2 "github.com/nocturna-ta/golib/response"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/interfaces/dao"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

func (m *Module) GetVoteResultByID(ctx context.Context, id string) (*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultByID")
	defer span.End()

	if id == "" {
		return nil, &custerr.ErrChain{
			Message: "vote result ID is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	result, err := m.voteResultRepo.GetVoteResultByID(ctx, id)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultByID] Failed to get vote result")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "vote result not found",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	return &response.VoteResultResponse{
		ID:              result.ID,
		VoterID:         result.VoterID,
		ElectionPairID:  result.ElectionPairID,
		Region:          result.Region,
		Status:          result.Status,
		TransactionHash: result.TransactionHash,
		ErrorMessage:    result.ErrorMessage,
		VotedAt:         result.VotedAt,
		ProcessedAt:     result.ProcessedAt,
		CreatedAt:       result.CreatedAt,
		UpdatedAt:       result.UpdatedAt,
	}, nil
}

func (m *Module) GetVoteResultsByElectionPair(ctx context.Context, electionPairID string, limit, offset int) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByElectionPair")
	defer span.End()

	if electionPairID == "" {
		return nil, &custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	results, err := m.voteResultRepo.GetVoteResultsByElectionPair(ctx, electionPairID, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"electionPairID": electionPairID,
			"limit":          limit,
			"offset":         offset,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByElectionPair] Failed to get vote results by election pair")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given election pair",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given election pair",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}
	return responses, nil
}

func (m *Module) GetVoteResultsByRegion(ctx context.Context, region string, limit, offset int) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByRegion")
	defer span.End()

	if region == "" {
		return nil, &custerr.ErrChain{
			Message: "region is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	results, err := m.voteResultRepo.GetVoteResultsByRegion(ctx, region, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByRegion] Failed to get vote results by region")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given region",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given region",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) GetVoteResultsByStatus(ctx context.Context, status string, limit, offset int) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByStatus")
	defer span.End()

	if status == "" {
		return nil, &custerr.ErrChain{
			Message: "status is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	results, err := m.voteResultRepo.GetVoteResultsByStatus(ctx, status, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"status": status,
			"limit":  limit,
			"offset": offset,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByStatus] Failed to get vote results by status")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given status",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}
	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}
	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given status",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) GetVoteResultsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByDateRange")
	defer span.End()

	if startDate.IsZero() || endDate.IsZero() {
		return nil, &custerr.ErrChain{
			Message: "start date and end date are required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if startDate.After(endDate) {
		return nil, &custerr.ErrChain{
			Message: "start date cannot be after end date",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if limit <= 0 {
		limit = 50
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	results, err := m.voteResultRepo.GetVoteResultsByDateRange(ctx, startDate, endDate, limit, offset)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"startDate": startDate,
			"endDate":   endDate,
			"limit":     limit,
			"offset":    offset,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByDateRange] Failed to get vote results by date range")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given date range",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given date range",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) GetElectionResults(ctx context.Context, electionPairID string) (*response.ElectionVoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetElectionResults")
	defer span.End()

	if electionPairID == "" {
		return nil, &custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	result, err := m.voteResultRepo.GetElectionResults(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"electionPairID": electionPairID,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetElectionResults] Failed to get election results")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "election results not found for the given election pair",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	return &response.ElectionVoteResultResponse{
		ElectionPairID: result.ElectionPairID,
		Region:         result.Region,
		TotalVotes:     result.TotalVotes,
		ConfirmedVotes: result.ConfirmedVotes,
		PendingVotes:   result.PendingVotes,
		ErrorVotes:     result.ErrorVotes,
		LastUpdated:    result.LastUpdated,
	}, nil
}

func (m *Module) GetElectionResultsByRegion(ctx context.Context, region string) ([]*response.ElectionVoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetElectionResultsByRegion")
	defer span.End()

	if region == "" {
		return nil, &custerr.ErrChain{
			Message: "region is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	results, err := m.voteResultRepo.GetElectionResultsByRegion(ctx, region)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetElectionResultsByRegion] Failed to get election results by region")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no election results found for the given region",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.ElectionVoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.ElectionVoteResultResponse{
			ElectionPairID: result.ElectionPairID,
			Region:         result.Region,
			TotalVotes:     result.TotalVotes,
			ConfirmedVotes: result.ConfirmedVotes,
			PendingVotes:   result.PendingVotes,
			ErrorVotes:     result.ErrorVotes,
			LastUpdated:    result.LastUpdated,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no election results found for the given region",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) GetRegionResults(ctx context.Context, region string) (*response.RegionVoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetRegionResults")
	defer span.End()

	if region == "" {
		return nil, &custerr.ErrChain{
			Message: "region is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	result, err := m.voteResultRepo.GetRegionResults(ctx, region)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetRegionResults] Failed to get region results")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "region results not found",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	return &response.RegionVoteResultResponse{
		Region:         result.Region,
		TotalVotes:     result.TotalVotes,
		ConfirmedVotes: result.ConfirmedVotes,
		PendingVotes:   result.PendingVotes,
		ErrorVotes:     result.ErrorVotes,
		LastUpdated:    result.LastUpdated,
	}, nil
}

func (m *Module) GetRegionStatistics(ctx context.Context) ([]*response.RegionVoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetRegionStatistics")
	defer span.End()

	results, err := m.voteResultRepo.GetRegionStatistics(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetRegionStatistics] Failed to get region statistics")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no region statistics found",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.RegionVoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.RegionVoteResultResponse{
			Region:         result.Region,
			TotalVotes:     result.TotalVotes,
			ConfirmedVotes: result.ConfirmedVotes,
			PendingVotes:   result.PendingVotes,
			ErrorVotes:     result.ErrorVotes,
			LastUpdated:    result.LastUpdated,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no region statistics found",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) GetOverallStatistics(ctx context.Context) (*response.VoteStatisticsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetOverallStatistics")
	defer span.End()

	stats, err := m.voteResultRepo.GetOverallStatistics(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetOverallStatistics] Failed to get overall statistics")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "overall statistics not found",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	return &response.VoteStatisticsResponse{
		TotalVotes:     stats.TotalVotes,
		ConfirmedVotes: stats.ConfirmedVotes,
		PendingVotes:   stats.PendingVotes,
		ErrorVotes:     stats.ErrorVotes,
		SuccessRate:    stats.SuccessRate,
		LastUpdated:    stats.LastUpdated,
	}, nil
}

func (m *Module) GetDailyStatistics(ctx context.Context, startDate, endDate time.Time) ([]*response.VoteStatisticsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetDailyStatistics")
	defer span.End()

	if startDate.IsZero() || endDate.IsZero() {
		return nil, &custerr.ErrChain{
			Message: "start date and end date are required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	if startDate.After(endDate) {
		return nil, &custerr.ErrChain{
			Message: "start date cannot be after end date",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	stats, err := m.voteResultRepo.GetDailyStatistics(ctx, startDate, endDate)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"startDate": startDate,
			"endDate":   endDate,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetDailyStatistics] Failed to get daily statistics")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no daily statistics found for the given date range",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.VoteStatisticsResponse
	for _, stat := range stats {
		responses = append(responses, &response.VoteStatisticsResponse{
			TotalVotes:     stat.TotalVotes,
			ConfirmedVotes: stat.ConfirmedVotes,
			PendingVotes:   stat.PendingVotes,
			ErrorVotes:     stat.ErrorVotes,
			SuccessRate:    stat.SuccessRate,
			LastUpdated:    stat.LastUpdated,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no daily statistics found for the given date range",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil
}

func (m *Module) CountVotesByStatus(ctx context.Context, status string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.CountVotesByStatus")
	defer span.End()

	if status == "" {
		return 0, &custerr.ErrChain{
			Message: "status is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	count, err := m.voteResultRepo.CountVotesByStatus(ctx, status)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"status": status,
		}).ErrorWithCtx(ctx, "[ResultUseCases.CountVotesByStatus] Failed to count votes by status")

		if errors.Is(err, dao.ErrNoResult) {
			return 0, &custerr.ErrChain{
				Message: "no votes found for the given status",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return 0, err
	}

	return count, nil
}

func (m *Module) CountVotesByElectionPair(ctx context.Context, electionPairID string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.CountVotesByElectionPair")
	defer span.End()

	if electionPairID == "" {
		return 0, &custerr.ErrChain{
			Message: "election pair ID is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	count, err := m.voteResultRepo.CountVotesByElectionPair(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":          err,
			"electionPairID": electionPairID,
		}).ErrorWithCtx(ctx, "[ResultUseCases.CountVotesByElectionPair] Failed to count votes by election pair")

		if errors.Is(err, dao.ErrNoResult) {
			return 0, &custerr.ErrChain{
				Message: "no votes found for the given election pair",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return 0, err
	}

	return count, nil
}

func (m *Module) CountVotesByRegion(ctx context.Context, region string) (uint64, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.CountVotesByRegion")
	defer span.End()

	if region == "" {
		return 0, &custerr.ErrChain{
			Message: "region is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	count, err := m.voteResultRepo.CountVotesByRegion(ctx, region)
	if err != nil {
		log.WithFields(log.Fields{
			"error":  err,
			"region": region,
		}).ErrorWithCtx(ctx, "[ResultUseCases.CountVotesByRegion] Failed to count votes by region")

		if errors.Is(err, dao.ErrNoResult) {
			return 0, &custerr.ErrChain{
				Message: "no votes found for the given region",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return 0, err
	}

	return count, nil
}

func (m *Module) GetVoteResultsByHour(ctx context.Context, date time.Time) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByHour")
	defer span.End()

	if date.IsZero() {
		return nil, &custerr.ErrChain{
			Message: "date is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	results, err := m.voteResultRepo.GetVoteResultsByHour(ctx, date)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"date":  date,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByHour] Failed to get vote results by hour")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given date",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}
	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}

	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given date",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}

	return responses, nil

}

func (m *Module) GetVoteResultsByDay(ctx context.Context, date time.Time) ([]*response.VoteResultResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "ResultUseCases.GetVoteResultsByDay")
	defer span.End()

	if date.IsZero() {
		return nil, &custerr.ErrChain{
			Message: "date is required",
			Code:    400,
			Type:    response2.ErrBadRequest,
		}
	}

	results, err := m.voteResultRepo.GetVoteResultsByDay(ctx, date)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"date":  date,
		}).ErrorWithCtx(ctx, "[ResultUseCases.GetVoteResultsByDay] Failed to get vote results by day")

		if errors.Is(err, dao.ErrNoResult) {
			return nil, &custerr.ErrChain{
				Message: "no vote results found for the given date",
				Code:    404,
				Type:    response2.ErrNotFound,
			}
		}
		return nil, err
	}

	var responses []*response.VoteResultResponse
	for _, result := range results {
		responses = append(responses, &response.VoteResultResponse{
			ID:              result.ID,
			VoterID:         result.VoterID,
			ElectionPairID:  result.ElectionPairID,
			Region:          result.Region,
			Status:          result.Status,
			TransactionHash: result.TransactionHash,
			ErrorMessage:    result.ErrorMessage,
			VotedAt:         result.VotedAt,
			ProcessedAt:     result.ProcessedAt,
			CreatedAt:       result.CreatedAt,
			UpdatedAt:       result.UpdatedAt,
		})
	}
	if len(responses) == 0 {
		return nil, &custerr.ErrChain{
			Message: "no vote results found for the given date",
			Code:    404,
			Type:    response2.ErrNotFound,
		}
	}
	return responses, nil
}
