package usecases

import (
	"context"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

type VoteResultUseCases interface {
	// Vote Result queries
	GetVoteResultByID(ctx context.Context, id string) (*response.VoteResultResponse, error)
	GetVoteResultsByElectionPair(ctx context.Context, electionPairID string, limit, offset int) ([]*response.VoteResultResponse, error)
	GetVoteResultsByRegion(ctx context.Context, region string, limit, offset int) ([]*response.VoteResultResponse, error)
	GetVoteResultsByStatus(ctx context.Context, status string, limit, offset int) ([]*response.VoteResultResponse, error)
	GetVoteResultsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*response.VoteResultResponse, error)

	// Election Results
	GetElectionResults(ctx context.Context, electionPairID string) (*response.ElectionVoteResultResponse, error)
	GetElectionResultsByRegion(ctx context.Context, region string) ([]*response.ElectionVoteResultResponse, error)

	// Region Results
	GetRegionResults(ctx context.Context, region string) (*response.RegionVoteResultResponse, error)
	GetRegionStatistics(ctx context.Context) ([]*response.RegionVoteResultResponse, error)

	// Statistics
	GetOverallStatistics(ctx context.Context) (*response.VoteStatisticsResponse, error)
	GetDailyStatistics(ctx context.Context, startDate, endDate time.Time) ([]*response.VoteStatisticsResponse, error)

	// Count operations
	CountVotesByStatus(ctx context.Context, status string) (uint64, error)
	CountVotesByElectionPair(ctx context.Context, electionPairID string) (uint64, error)
	CountVotesByRegion(ctx context.Context, region string) (uint64, error)

	// Time-based queries
	GetVoteResultsByHour(ctx context.Context, date time.Time) ([]*response.VoteResultResponse, error)
	GetVoteResultsByDay(ctx context.Context, date time.Time) ([]*response.VoteResultResponse, error)
}
