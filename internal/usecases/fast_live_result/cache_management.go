package fast_live_result

import (
	"context"
	"errors"
	"fmt"
	"github.com/nocturna-ta/golib/cache"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/golib/tracing"
	"github.com/nocturna-ta/result/internal/infrastructures/websocket"
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

	// Create WebSocket message using their existing structure
	message := &websocket.LiveMessage{
		Type:      "fast_live_results_update", // New message type for fast results
		Timestamp: time.Now(),
		Data:      electionResults,
		Filter: &websocket.MessageFilter{
			ElectionPairID: electionPairID,
		},
	}

	select {
	case m.wsHub.Broadcast <- message:
	default:
		log.WarnWithCtx(ctx, "[FastLiveResultUseCase] Broadcast channel full, dropping live results update")
	}

	log.WithFields(log.Fields{
		"election_id": electionPairID,
		"clients":     m.wsHub.GetClientCount(),
		"regions":     len(electionResults.RegionResults),
		"cities":      len(electionResults.TopCities),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Live results broadcasted")

	return nil
}

func (m *Module) BroadcastIncrementalUpdate(ctx context.Context, updateData *response.IncrementalUpdateData) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCase.BroadcastIncrementalUpdate")
	defer span.End()

	if m.wsHub.GetClientCount() == 0 {
		return nil // No clients connected, skip broadcast
	}

	// Create incremental update message using their existing structure
	message := &websocket.LiveMessage{
		Type:      "fast_incremental_update", // New message type for incremental updates
		Timestamp: updateData.Timestamp,
		Data:      updateData,
		Filter: &websocket.MessageFilter{
			ElectionPairID: updateData.ElectionID,
			Region:         updateData.Region,
		},
	}

	// Send to broadcast channel (their existing pattern)
	select {
	case m.wsHub.Broadcast <- message:
	default:
		log.WarnWithCtx(ctx, "[FastLiveResultUseCase] Broadcast channel full, dropping incremental update")
	}

	log.WithFields(log.Fields{
		"election_id": updateData.ElectionID,
		"region":      updateData.Region,
		"city_name":   updateData.CityName,
		"new_votes":   updateData.NewVotes,
		"clients":     m.wsHub.GetClientCount(),
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Incremental update broadcasted")

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
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase.InvalidateElectionCache] Failed to invalidate cache key")
		}

		log.WithFields(log.Fields{
			"election_id":      electionPairID,
			"invalidated_keys": len(keys),
		}).InfoWithCtx(ctx, "[FastLiveResultUseCase.InvalidateElectionCache] Cache invalidated")
	}
	return nil
}

func (m *Module) InvalidateCityCache(ctx context.Context, cityName string) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.InvalidateCityCache")
	defer span.End()

	key := fmt.Sprintf(CacheKeyCityResults, cityName)

	if err := m.redisCache.Delete(ctx, key); err != nil && !errors.Is(err, cache.ErrNotFound) {
		log.WithFields(log.Fields{
			"error":     err,
			"city_name": cityName,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase.InvalidateCityCache] Failed to invalidate city cache key")
		return err
	}
	return nil
}

func (m *Module) InvalidateAllCache(ctx context.Context) error {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.InvalidateAllCache")
	defer span.End()

	patterns := []string{
		"fast_live:election:*",
		"fast_live:city:*",
		"fast_live:summary:*",
		"fast_live:rankings:*",
	}

	// Also delete specific single keys
	singleKeys := []string{
		CacheKeyAllElections,
		CacheKeyStatistics,
	}

	var totalKeys int

	// Delete pattern-based keys
	for _, pattern := range patterns {
		keys := m.redisCache.GetKeys(ctx, pattern)
		totalKeys += len(keys)

		// Delete each key individually since golib cache doesn't have pattern delete
		for _, key := range keys {
			// Remove namespace prefix if present
			cleanKey := key
			if strings.HasPrefix(key, "fast_live:") {
				cleanKey = strings.TrimPrefix(key, "fast_live:")
			}

			if err := m.redisCache.Delete(ctx, cleanKey); err != nil && err != cache.ErrNotFound {
				log.WithFields(log.Fields{
					"error": err,
					"key":   cleanKey,
				}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to invalidate cache key")
			}
		}
	}

	// Delete single keys
	for _, key := range singleKeys {
		if err := m.redisCache.Delete(ctx, key); err != nil && err != cache.ErrNotFound {
			log.WithFields(log.Fields{
				"error": err,
				"key":   key,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to invalidate single cache key")
		}
		totalKeys++
	}

	log.WithFields(log.Fields{
		"invalidated_keys": totalKeys,
	}).InfoWithCtx(ctx, "[FastLiveResultUseCase] All cache invalidated")

	// Reset statistics
	atomic.StoreUint64(&m.cacheStatistics.TotalHits, 0)
	atomic.StoreUint64(&m.cacheStatistics.TotalMisses, 0)

	return nil
}

func (m *Module) StartCacheWarming(ctx context.Context, activeElections []string) {
	go func() {
		ticker := time.NewTicker(15 * time.Second) // Warm cache every 15 seconds
		defer ticker.Stop()

		log.WithFields(log.Fields{
			"elections": len(activeElections),
			"interval":  "15s",
		}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Starting cache warming service")

		for {
			select {
			case <-ctx.Done():
				log.InfoWithCtx(ctx, "[FastLiveResultUseCase] Cache warming service stopped")
				return
			case <-ticker.C:
				m.warmCache(context.Background(), activeElections)
			}
		}
	}()
}

func (m *Module) StartIncrementalBroadcast(ctx context.Context, interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		log.WithFields(log.Fields{
			"interval": interval,
		}).InfoWithCtx(ctx, "[FastLiveResultUseCase] Starting incremental broadcast service")

		for {
			select {
			case <-ctx.Done():
				log.InfoWithCtx(ctx, "[FastLiveResultUseCase] Incremental broadcast service stopped")
				return
			case <-ticker.C:
				if m.wsHub.GetClientCount() > 0 {
					// Get active elections from cache keys
					electionKeys := m.redisCache.GetKeys(context.Background(), "fast_live:election:*")
					for _, key := range electionKeys {
						parts := strings.Split(key, ":")
						if len(parts) >= 3 {
							electionID := parts[2]
							// Broadcast updates for each active election
							m.BroadcastLiveResultsUpdate(context.Background(), electionID)
						}
					}
				}
			}
		}
	}()
}

func (m *Module) GetCacheStatistics(ctx context.Context) (*response.CacheStatisticsResponse, error) {
	span, ctx := tracing.StartSpanFromContext(ctx, "FastLiveResultUseCases.GetCacheStatistics")
	defer span.End()

	totalHits := atomic.LoadUint64(&m.cacheStatistics.TotalHits)
	totalMisses := atomic.LoadUint64(&m.cacheStatistics.TotalMisses)

	var hitRate float64
	if totalHits+totalMisses > 0 {
		hitRate = float64(totalHits) / float64(totalHits+totalMisses) * 100
	}

	// Get active elections by scanning cache keys
	electionKeys := m.redisCache.GetKeys(ctx, "fast_live:election:*")
	activeElections := make([]string, 0, len(electionKeys))
	for _, key := range electionKeys {
		// Extract election ID from key: fast_live:election:{election_id}
		parts := strings.Split(key, ":")
		if len(parts) >= 3 {
			activeElections = append(activeElections, parts[2])
		}
	}

	// Get Redis info
	redisInfo := make(map[string]interface{})
	// Note: This would need to be implemented based on your cache interface
	// For now, we'll provide basic info
	redisInfo["status"] = "connected"
	redisInfo["cache_keys"] = len(electionKeys) + len(m.redisCache.GetKeys(ctx, "fast_live:*"))

	return &response.CacheStatisticsResponse{
		HitRate:          hitRate,
		TotalHits:        totalHits,
		TotalMisses:      totalMisses,
		CacheSize:        "N/A",
		ConnectedClients: m.GetConnectedClientsCount(ctx),
		LastCacheRefresh: time.Now(),
		ActiveElections:  activeElections,
		RedisInfo:        redisInfo,
	}, nil
}

func (m *Module) GetConnectedClientsCount(ctx context.Context) int {
	return m.wsHub.GetClientCount()
}
