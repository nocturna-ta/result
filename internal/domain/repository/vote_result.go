package repository

import (
	"context"
	"github.com/nocturna-ta/result/internal/domain/model"
	"time"
)

type VoteResultRepository interface {
	// Insert operations
	InsertVoteResult(ctx context.Context, result *model.VoteResult) error
	UpdateVoteResult(ctx context.Context, result *model.VoteResult) error

	// Read operations
	GetVoteResultByID(ctx context.Context, id string) (*model.VoteResult, error)
	GetVoteResultsByElectionPair(ctx context.Context, electionPairID string, limit, offset int) ([]*model.VoteResult, error)
	GetVoteResultsByRegion(ctx context.Context, region string, limit, offset int) ([]*model.VoteResult, error)
	GetVoteResultsByStatus(ctx context.Context, status string, limit, offset int) ([]*model.VoteResult, error)

	// Statistics operations
	GetElectionResults(ctx context.Context, electionPairID string) (*model.ElectionResult, error)
	GetRegionResults(ctx context.Context, region string) (*model.RegionResult, error)
	GetOverallStatistics(ctx context.Context) (*model.VoteStatistics, error)

	// Advanced queries
	GetVoteResultsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*model.VoteResult, error)
	GetElectionResultsByRegion(ctx context.Context, region string) ([]*model.ElectionResult, error)
	GetRegionStatistics(ctx context.Context) ([]*model.RegionResult, error)

	// Count operations
	CountVotesByStatus(ctx context.Context, status string) (uint64, error)
	CountVotesByElectionPair(ctx context.Context, electionPairID string) (uint64, error)
	CountVotesByRegion(ctx context.Context, region string) (uint64, error)

	// Time-based aggregations
	GetVoteResultsByHour(ctx context.Context, date time.Time) ([]*model.VoteResult, error)
	GetVoteResultsByDay(ctx context.Context, date time.Time) ([]*model.VoteResult, error)
	GetDailyStatistics(ctx context.Context, startDate, endDate time.Time) ([]*model.VoteStatistics, error)
}
