package fast_live_result

import (
	"context"
	"errors"
	"fmt"
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"sync/atomic"
	"time"
)

const (
	CacheKeyElectionResults = "fast_live:election:%s"
	CacheKeyCityResults     = "fast_live:city:%s"
	CacheKeyElectionSummary = "fast_live:summary:%s"
	CacheKeyCityRankings    = "fast_live:rankings:%s:%d"
	CacheKeyAllElections    = "fast_live:all_elections"
	CacheKeyStatistics      = "fast_live:stats"

	DefaultCacheTTL      = 30
	SummaryCacheTTL      = 15
	AllElectionsCacheTTL = 20
	StatisticsCacheTTL   = 10
)

func (m *Module) GetLiveElectionResultsWithCache(ctx context.Context, electionPairID string) (*response.FastElectionResultsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetLiveElectionResultsWithCache")
	defer span.End()

	cacheKey := fmt.Sprintf(CacheKeyElectionResults, electionPairID)

	var cachedResult response.FastElectionResultsResponse
	err := m.redisCache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		atomic.AddUint64(&m.cacheStatistics.TotalHits, 1)
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"source":           "cache",
		}).DebugWithCtx(ctx, "[GetLiveElectionResultsWithCache]  Served election results from cache")
		return &cachedResult, nil
	}

	if !errors.Is(err, cache.ErrNotFound) {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"error":            err,
		}).WarnWithCtx(ctx, "[GetLiveElectionResultsWithCache] Failed to get election results from cache")
	}

	atomic.AddUint64(&m.cacheStatistics.TotalMisses, 1)

	liveResults, err := m.liveResultRepo.GetLiveElectionResults(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"error":            err,
		}).ErrorWithCtx(ctx, "[GetLiveElectionResultsWithCache] Failed to get live election results")
		return nil, err
	}

	summary, err := m.liveResultRepo.GetElectionSummary(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"error":            err,
		}).ErrorWithCtx(ctx, "[GetLiveElectionResultsWithCache] Failed to get election summary")
	}

	cityRankings, err := m.liveResultRepo.GetCityRankings(ctx, electionPairID, 10)
	if err != nil {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"error":            err,
		}).ErrorWithCtx(ctx, "[GetLiveElectionResultsWithCache] Failed to get city rankings")
	}

	response2 := m.buildElectionResultsResponse(liveResults, summary, cityRankings, electionPairID)

	go m.cacheElectionResults(context.Background(), cacheKey, response2)

	log.WithFields(log.Fields{
		"election_id": electionPairID,
		"source":      "database",
		"regions":     len(response2.RegionResults),
		"cities":      len(response2.TopCities),
	}).DebugWithCtx(ctx, "[FastLiveResultUseCase] Served election results from database")

	return response2, nil
}

func (m *Module) GetLiveCityResultsWithCache(ctx context.Context, cityName string) (*response.FastCityResultsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetLiveCityResultsWithCache")
	defer span.End()

	cacheKey := fmt.Sprintf(CacheKeyCityResults, cityName)

	var cachedResult response.FastCityResultsResponse
	err := m.redisCache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		atomic.AddUint64(&m.cacheStatistics.TotalHits, 1)
		log.WithFields(log.Fields{
			"city_name": cityName,
			"source":    "cache",
		}).DebugWithCtx(ctx, "[GetLiveCityResultsWithCache] Served city results from cache")
		return &cachedResult, nil
	}

	atomic.AddUint64(&m.cacheStatistics.TotalMisses, 1)

	cityResults, err := m.liveResultRepo.GetLiveCityResults(ctx, cityName)
	if err != nil {
		log.WithFields(log.Fields{
			"city_name": cityName,
			"error":     err,
		}).ErrorWithCtx(ctx, "[GetLiveCityResultsWithCache] Failed to get live city results")
		return nil, err
	}

	if len(cityResults) == 0 {
		return &response.FastCityResultsResponse{
			CityName:    cityName,
			LastUpdated: time.Now(),
		}, nil
	}

	response2 := m.buildCityResultsResponse(cityResults, cityName)

	go m.cacheCityResults(context.Background(), cacheKey, response2)

	return response2, nil
}

func (m *Module) GetElectionSummaryWithCache(ctx context.Context, electionPairID string) (*response.FastElectionSummaryResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetElectionSummaryWithCache")
	defer span.End()

	cacheKey := fmt.Sprintf(CacheKeyElectionSummary, electionPairID)

	var cachedResult response.FastElectionSummaryResponse
	err := m.redisCache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		atomic.AddUint64(&m.cacheStatistics.TotalHits, 1)
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"source":           "cache",
		}).DebugWithCtx(ctx, "[GetElectionSummaryWithCache] Served election summary from cache")
		return &cachedResult, nil
	}

	atomic.AddUint64(&m.cacheStatistics.TotalMisses, 1)

	summary, err := m.liveResultRepo.GetElectionSummary(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"error":            err,
		}).ErrorWithCtx(ctx, "[GetElectionSummaryWithCache] Failed to get election summary")
	}

	response2 := &response.FastElectionSummaryResponse{
		ElectionID:          summary.ElectionPairID,
		TotalUniqueVoters:   summary.TotalUniqueVoters,
		TotalRegions:        summary.TotalRegions,
		TotalConfirmedVotes: summary.TotalConfirmedVotes,
		OverallSuccessRate:  summary.OverallSuccessRate,
		LastUpdated:         summary.LastUpdated,
	}

	go m.cacheElectionSummary(context.Background(), cacheKey, response2)

	return response2, nil
}

func (m *Module) GetCityRankingsWithCache(ctx context.Context, electionPairID string, limit int) (*response.FastCityRankingsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetCityRankingsWithCache")
	defer span.End()

	if limit <= 0 || limit > 100 {
		limit = 20
	}

	cacheKey := fmt.Sprintf(CacheKeyCityRankings, electionPairID, limit)
	var cachedResult response.FastCityRankingsResponse
	err := m.redisCache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		atomic.AddUint64(&m.cacheStatistics.TotalHits, 1)
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"limit":            limit,
			"source":           "cache",
		}).DebugWithCtx(ctx, "[GetCityRankingsWithCache] Served city rankings from cache")
		return &cachedResult, nil
	}

	atomic.AddUint64(&m.cacheStatistics.TotalMisses, 1)

	rankings, err := m.liveResultRepo.GetCityRankings(ctx, electionPairID, limit)
	if err != nil {
		log.WithFields(log.Fields{
			"election_pair_id": electionPairID,
			"limit":            limit,
			"error":            err,
		}).ErrorWithCtx(ctx, "[GetCityRankingsWithCache] Failed to get city rankings")
		return nil, err
	}

	response2 := m.buildCityRankingsResponse(rankings, electionPairID)

	go m.cacheCityRankings(context.Background(), cacheKey, response2)

	return response2, nil
}

func (m *Module) GetAllElectionsSummaryWithCache(ctx context.Context) (*response.FastAllElectionsSummaryResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetAllElectionsSummaryWithCache")
	defer span.End()

	cacheKey := CacheKeyAllElections

	var cachedResult response.FastAllElectionsSummaryResponse
	err := m.redisCache.GetObject(ctx, cacheKey, &cachedResult)
	if err == nil {
		atomic.AddUint64(&m.cacheStatistics.TotalHits, 1)
		log.WithFields(log.Fields{
			"source": "cache",
		}).DebugWithCtx(ctx, "[GetAllElectionsSummaryWithCache] Served all elections summary from cache")
		return &cachedResult, nil
	}

	atomic.AddUint64(&m.cacheStatistics.TotalMisses, 1)

	summaries, err := m.liveResultRepo.GetAllElectionsSummary(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[GetAllElectionsSummaryWithCache] Failed to get all elections summary")
		return nil, err
	}

	response2 := &response.FastAllElectionsSummaryResponse{
		Elections:      make([]response.FastElectionSummaryResponse, 0, len(summaries)),
		LastUpdated:    time.Now(),
		TotalElections: uint64(len(summaries)),
	}

	for _, summary := range summaries {
		response2.Elections = append(response2.Elections, response.FastElectionSummaryResponse{
			ElectionID:          summary.ElectionPairID,
			TotalUniqueVoters:   summary.TotalUniqueVoters,
			TotalRegions:        summary.TotalRegions,
			TotalConfirmedVotes: summary.TotalConfirmedVotes,
			OverallSuccessRate:  summary.OverallSuccessRate,
			LastUpdated:         summary.LastUpdated,
		})
	}

	go m.cacheAllElections(context.Background(), cacheKey, response2)

	return response2, nil
}
