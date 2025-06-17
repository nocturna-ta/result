package fast_live_result

import (
	"context"
	"errors"
	"fmt"
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"strings"
	"sync/atomic"
	"time"
)

func (m *Module) BroadcastLiveResultsUpdate(ctx context.Context, electionPairID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.BroadcastLiveResultsUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	// Get fresh data from cache/database
	electionResults, err := m.GetLiveElectionResultsWithCache(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"election_id": electionPairID,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get election results for broadcast")
		return err
	}

	m.wsHub.BroadcastFastLiveResults(electionResults)

	log.WithFields(log.Fields{
		"election_id": electionPairID,
		"clients":     m.wsHub.GetClientCount(),
		"regions":     len(electionResults.RegionResults),
		"cities":      len(electionResults.TopCities),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast live results broadcasted")

	return nil
}

func (m *Module) BroadcastIncrementalUpdate(ctx context.Context, updateData *response.IncrementalUpdateData) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastIncrementalUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	log.WithFields(log.Fields{
		"election_id": updateData.ElectionID,
		"region":      updateData.Region,
		"city_name":   updateData.CityName,
		"new_votes":   updateData.NewVotes,
		"clients":     m.wsHub.GetClientCount(),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast incremental update broadcasted")

	return nil
}

func (m *Module) BroadcastCityResultsUpdate(ctx context.Context, cityName string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastCityResultsUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	cityResults, err := m.GetLiveCityResultsWithCache(ctx, cityName)
	if err != nil {
		log.WithFields(log.Fields{
			"error":     err,
			"city_name": cityName,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get city results for broadcast")
		return err
	}

	m.wsHub.BroadcastFastCityResults(cityResults)

	log.WithFields(log.Fields{
		"city_name":      cityName,
		"clients":        m.wsHub.GetClientCount(),
		"total_votes":    cityResults.TotalVotes,
		"total_voters":   cityResults.TotalVoters,
		"election_count": len(cityResults.ElectionResults),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast city results broadcasted")

	return nil
}

func (m *Module) BroadcastRankingsUpdate(ctx context.Context, electionPairID string, limit int) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastRankingsUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	// Get fresh rankings data from cache/database
	rankings, err := m.GetCityRankingsWithCache(ctx, electionPairID, limit)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"election_id": electionPairID,
			"limit":       limit,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get rankings for broadcast")
		return err
	}

	// Use the Hub's dedicated method for rankings
	m.wsHub.BroadcastFastRankings(rankings)

	log.WithFields(log.Fields{
		"election_id":  electionPairID,
		"clients":      m.wsHub.GetClientCount(),
		"total_cities": rankings.TotalCities,
		"ranking_size": len(rankings.Rankings),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast rankings broadcasted")

	return nil
}

func (m *Module) BroadcastElectionSummaryUpdate(ctx context.Context, electionPairID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastElectionSummaryUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	// Get fresh summary data from cache/database
	summary, err := m.GetElectionSummaryWithCache(ctx, electionPairID)
	if err != nil {
		log.WithFields(log.Fields{
			"error":       err,
			"election_id": electionPairID,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get election summary for broadcast")
		return err
	}

	// Use the Hub's dedicated method for election summary
	m.wsHub.BroadcastFastElectionSummary(summary)

	log.WithFields(log.Fields{
		"election_id":           electionPairID,
		"clients":               m.wsHub.GetClientCount(),
		"total_unique_voters":   summary.TotalUniqueVoters,
		"total_confirmed_votes": summary.TotalConfirmedVotes,
		"overall_success_rate":  summary.OverallSuccessRate,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast election summary broadcasted")

	return nil
}

func (m *Module) BroadcastAllElectionsUpdate(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastAllElectionsUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	// Get fresh all elections data from cache/database
	allElections, err := m.GetAllElectionsSummaryWithCache(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get all elections summary for broadcast")
		return err
	}

	// Use the Hub's dedicated method for all elections
	m.wsHub.BroadcastFastAllElections(allElections)

	log.WithFields(log.Fields{
		"clients":         m.wsHub.GetClientCount(),
		"total_elections": allElections.TotalElections,
		"election_count":  len(allElections.Elections),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Fast all elections summary broadcasted")

	return nil
}

func (m *Module) BroadcastBulkUpdate(ctx context.Context, electionPairID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastBulkUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	var errors []error

	// Broadcast live results update
	if err := m.BroadcastLiveResultsUpdate(ctx, electionPairID); err != nil {
		errors = append(errors, fmt.Errorf("live results: %w", err))
	}

	// Broadcast election summary
	if err := m.BroadcastElectionSummaryUpdate(ctx, electionPairID); err != nil {
		errors = append(errors, fmt.Errorf("election summary: %w", err))
	}

	// Broadcast rankings (top 100)
	if err := m.BroadcastRankingsUpdate(ctx, electionPairID, 100); err != nil {
		errors = append(errors, fmt.Errorf("rankings: %w", err))
	}

	// Broadcast all elections summary
	if err := m.BroadcastAllElectionsUpdate(ctx); err != nil {
		errors = append(errors, fmt.Errorf("all elections: %w", err))
	}

	if len(errors) > 0 {
		var errorMessages []string
		for _, err := range errors {
			errorMessages = append(errorMessages, err.Error())
		}
		return fmt.Errorf("bulk broadcast errors: %s", strings.Join(errorMessages, "; "))
	}

	log.WithFields(log.Fields{
		"election_id": electionPairID,
		"clients":     m.wsHub.GetClientCount(),
		"broadcasts":  4,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Bulk broadcast completed successfully")

	return nil
}

func (m *Module) InvalidateElectionCache(ctx context.Context, electionPairID string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.InvalidateElectionCache")
	defer span.End()

	keys := []string{
		fmt.Sprintf(CacheKeyElectionResults, electionPairID),
		fmt.Sprintf(CacheKeyElectionSummary, electionPairID),
		CacheKeyAllElections,
		CacheKeyStatistics,
	}

	rankingPattern := fmt.Sprintf("fast_live:rankings:%s:*", electionPairID)
	rankingKeys := m.redisCache.GetKeys(ctx, rankingPattern)
	keys = append(keys, rankingKeys...)

	for _, key := range keys {
		if err := m.redisCache.Delete(ctx, key); err != nil && !errors.Is(err, cache.ErrNotFound) {
			log.WithFields(log.Fields{
				"error":       err,
				"key":         key,
				"election_id": electionPairID,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to delete cache key")
		}
	}

	log.WithFields(log.Fields{
		"election_id":  electionPairID,
		"keys_deleted": len(keys),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Election cache invalidated")

	return nil
}

func (m *Module) InvalidateCityCache(ctx context.Context, cityName string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.InvalidateCityCache")
	defer span.End()

	key := fmt.Sprintf(CacheKeyCityResults, cityName)
	if err := m.redisCache.Delete(ctx, key); err != nil && !errors.Is(err, cache.ErrNotFound) {
		log.WithFields(log.Fields{
			"error":     err,
			"key":       key,
			"city_name": cityName,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to delete city cache key")
		return err
	}

	log.WithFields(log.Fields{
		"city_name": cityName,
		"key":       key,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] City cache invalidated")

	return nil
}

func (m *Module) InvalidateAllCache(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.InvalidateAllCache")
	defer span.End()

	// Get all fast live result cache keys
	pattern := "fast_live:*"
	keys := m.redisCache.GetKeys(ctx, pattern)

	deletedCount := 0
	for _, key := range keys {
		if err := m.redisCache.Delete(ctx, key); err != nil && !errors.Is(err, cache.ErrNotFound) {
			log.WithFields(log.Fields{
				"error": err,
				"key":   key,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to delete cache key during bulk invalidation")
		} else {
			deletedCount++
		}
	}

	log.WithFields(log.Fields{
		"keys_found":    len(keys),
		"keys_deleted":  deletedCount,
		"cache_pattern": pattern,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] All cache invalidated")

	return nil
}

// Cache warming and background tasks
func (m *Module) StartCacheWarming(ctx context.Context, activeElections []string) {
	go func() {
		span, warmCtx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.CacheWarming")
		defer span.End()

		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for {
			select {
			case <-warmCtx.Done():
				log.InfoWithCtx(warmCtx, "[FastLiveResultUseCase] Cache warming stopped")
				return
			case <-ticker.C:
				m.warmElectionCaches(warmCtx, activeElections)
			}
		}
	}()
}

func (m *Module) warmElectionCaches(ctx context.Context, electionIDs []string) {
	for _, electionID := range electionIDs {
		// Warm election results cache
		if _, err := m.GetLiveElectionResultsWithCache(ctx, electionID); err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"election_id": electionID,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm election results cache")
		}

		// Warm election summary cache
		if _, err := m.GetElectionSummaryWithCache(ctx, electionID); err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"election_id": electionID,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm election summary cache")
		}

		// Warm rankings cache
		if _, err := m.GetCityRankingsWithCache(ctx, electionID, 100); err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"election_id": electionID,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm rankings cache")
		}
	}

	// Warm all elections summary
	if _, err := m.GetAllElectionsSummaryWithCache(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm all elections cache")
	}

	log.WithFields(log.Fields{
		"elections_warmed": len(electionIDs),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Cache warming cycle completed")
}

func (m *Module) StartIncrementalBroadcast(ctx context.Context, interval time.Duration) {
	go func() {
		span, broadcastCtx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.IncrementalBroadcast")
		defer span.End()

		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-broadcastCtx.Done():
				log.InfoWithCtx(broadcastCtx, "[FastLiveResultUseCase] Incremental broadcast stopped")
				return
			case <-ticker.C:
				if m.wsHub.GetClientCount() > 0 {
					m.performIncrementalBroadcast(broadcastCtx)
				}
			}
		}
	}()
}

func (m *Module) performIncrementalBroadcast(ctx context.Context) {
	// Get all active elections
	allElections, err := m.GetAllElectionsSummaryWithCache(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).ErrorWithCtx(ctx, "[FastLiveResultUseCase] Failed to get elections for incremental broadcast")
		return
	}

	broadcastCount := int64(0)
	for _, election := range allElections.Elections {
		// Broadcast live results for each active election
		if err := m.BroadcastLiveResultsUpdate(ctx, election.ElectionID); err != nil {
			log.WithFields(log.Fields{
				"error":       err,
				"election_id": election.ElectionID,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed incremental broadcast for election")
		} else {
			atomic.AddInt64(&broadcastCount, 1)
		}
	}

	// Broadcast all elections summary
	if err := m.BroadcastAllElectionsUpdate(ctx); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to broadcast all elections summary")
	} else {
		atomic.AddInt64(&broadcastCount, 1)
	}

	log.WithFields(log.Fields{
		"clients":         m.wsHub.GetClientCount(),
		"elections":       len(allElections.Elections),
		"broadcasts_sent": broadcastCount,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Incremental broadcast cycle completed")
}

func (m *Module) GetCacheStatistics(ctx context.Context) (*response.CacheStatisticsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m *Module) GetConnectedClientsCount(ctx context.Context) int {
	//TODO implement me
	panic("implement me")
}
