package usecases

import (
	"context"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"time"
)

type FastLiveResultUseCases interface {
	GetLiveElectionResultsWithCache(ctx context.Context, electionPairID string) (*response.FastElectionResultsResponse, error)
	GetLiveCityResultsWithCache(ctx context.Context, cityName string) (*response.FastCityResultsResponse, error)
	GetElectionSummaryWithCache(ctx context.Context, electionPairID string) (*response.FastElectionSummaryResponse, error)
	GetCityRankingsWithCache(ctx context.Context, electionPairID string, limit int) (*response.FastCityRankingsResponse, error)
	GetAllElectionsSummaryWithCache(ctx context.Context) (*response.FastAllElectionsSummaryResponse, error)

	BroadcastLiveResultsUpdate(ctx context.Context, electionPairID string) error
	BroadcastIncrementalUpdate(ctx context.Context, updateData *response.IncrementalUpdateData) error

	InvalidateElectionCache(ctx context.Context, electionPairID string) error
	InvalidateCityCache(ctx context.Context, cityName string) error
	InvalidateAllCache(ctx context.Context) error

	StartCacheWarming(ctx context.Context, activeElections []string)
	StartIncrementalBroadcast(ctx context.Context, interval time.Duration)

	GetCacheStatistics(ctx context.Context) (*response.CacheStatisticsResponse, error)
	GetConnectedClientsCount(ctx context.Context) int
}
