package usecases

import (
	"context"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

type LiveResultUseCases interface {
	GetLiveElectionResultsWithCache(ctx context.Context, electionPairID string) (*response.ElectionResultsResponse, error)
	GetLiveCityResultsWithCache(ctx context.Context, cityName string) (*response.CityResultsResponse, error)
	GetElectionSummaryWithCache(ctx context.Context, electionPairID string) (*response.ElectionSummaryResponse, error)
	GetCityRankingsWithCache(ctx context.Context, electionPairID string, limit int) (*response.CityRankingsResponse, error)
	GetAllElectionsSummaryWithCache(ctx context.Context) (*response.AllElectionsSummaryResponse, error)

	// Individual broadcast methods
	BroadcastVoteUpdate(ctx context.Context, voteID string) error
	BroadcastLiveResultsUpdate(ctx context.Context, electionPairID string) error
	BroadcastIncrementalUpdate(ctx context.Context, updateData *response.IncrementalUpdateData) error
	BroadcastCityResultsUpdate(ctx context.Context, cityName string) error
	BroadcastRankingsUpdate(ctx context.Context, electionPairID string, limit int) error
	BroadcastElectionSummaryUpdate(ctx context.Context, electionPairID string) error
	BroadcastAllElectionsUpdate(ctx context.Context) error

	// Bulk broadcast method
	BroadcastBulkUpdate(ctx context.Context, electionPairID string) error

	// Cache management methods
	InvalidateElectionCache(ctx context.Context, electionPairID string) error
	InvalidateCityCache(ctx context.Context, cityName string) error
	InvalidateAllCache(ctx context.Context) error

	// Background task methods
	StartCacheWarming(ctx context.Context, activeElections []string)
	StartIncrementalBroadcast(ctx context.Context, interval time.Duration)

	// Statistics and monitoring
	GetConnectedClientsCount(ctx context.Context) int
}
