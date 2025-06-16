package fast_live_result

import (
	"context"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"math"
	"time"
)

func (m *Module) buildElectionResultsResponse(liveResults []*model.LiveElectionResult, summary *model.ElectionSummary, cityRankings []*model.CityRanking, electionPairID string) *response.FastElectionResultsResponse {
	response2 := &response.FastElectionResultsResponse{
		ElectionID:    electionPairID,
		RegionResults: make([]response.RegionResultSummary, 0, len(liveResults)),
		TopCities:     make([]response.CityResultSummary, 0, len(cityRankings)),
	}

	// Add region results
	for _, result := range liveResults {
		response2.RegionResults = append(response2.RegionResults, response.RegionResultSummary{
			Region:     result.Region,
			Votes:      result.ConfirmedVotes,
			Percentage: result.SuccessPercentage,
		})
	}

	// Add summary data if available
	if summary != nil {
		response2.TotalVotes = summary.TotalConfirmedVotes
		response2.TotalVoters = summary.TotalUniqueVoters
		response2.LastUpdated = summary.LastUpdated
		response2.OverallPercentage = summary.OverallSuccessRate
		response2.OverallStats = response.ElectionStatsSummary{
			TotalVoters:    summary.TotalUniqueVoters,
			TotalRegions:   summary.TotalRegions,
			SuccessRate:    summary.OverallSuccessRate,
			ActiveRegions:  uint64(len(liveResults)),
			CompletionRate: calculateCompletionRate(summary),
		}

		// Calculate votes per second
		if !summary.LastUpdated.IsZero() {
			timeSinceStart := time.Since(summary.LastUpdated.Add(-24 * time.Hour))
			if timeSinceStart.Seconds() > 0 {
				response2.OverallStats.VotesPerSecond = float64(summary.TotalConfirmedVotes) / timeSinceStart.Seconds()
			}
		}
	}

	// Add top cities
	for _, city := range cityRankings {
		response2.TopCities = append(response2.TopCities, response.CityResultSummary{
			City:       city.CityName,
			Votes:      city.ConfirmedVotes,
			Voters:     city.UniqueVoters,
			Percentage: city.ParticipationRate,
			Rank:       city.CityRank,
		})
	}

	return response2
}

func (m *Module) buildCityResultsResponse(cityResults []*model.LiveCityResult, cityName string) *response.FastCityResultsResponse {
	response2 := &response.FastCityResultsResponse{
		CityName:        cityName,
		ElectionResults: make([]response.CityElectionResult, 0, len(cityResults)),
		LastUpdated:     time.Now(),
	}

	var totalVoters, totalVotes uint64
	var leadingElection string
	var maxVotes uint64

	for _, result := range cityResults {
		response2.ElectionResults = append(response2.ElectionResults, response.CityElectionResult{
			ElectionPairID:  result.ElectionPairID,
			Votes:           result.ConfirmedVotes,
			Percentage:      result.CandidatePercentageInCity,
			VoteSuccessRate: result.VoteSuccessRate,
		})

		totalVoters += result.TotalUniqueVoters
		totalVotes += result.ConfirmedVotes

		if result.ConfirmedVotes > maxVotes {
			maxVotes = result.ConfirmedVotes
			leadingElection = result.ElectionPairID
		}

		if response2.LastUpdated.Before(result.LastUpdated) {
			response2.LastUpdated = result.LastUpdated
		}
	}

	response2.TotalVoters = totalVoters
	response2.TotalVotes = totalVotes

	// Calculate city statistics
	response2.CityStats = response.CityStatistics{
		TotalElections:          uint64(len(cityResults)),
		LeadingElection:         leadingElection,
		ParticipationRate:       calculateParticipationRate(totalVotes, totalVoters),
		AverageVotesPerElection: float64(totalVotes) / float64(len(cityResults)),
	}

	return response2
}

func (m *Module) buildCityRankingsResponse(rankings []*model.CityRanking, electionPairID string) *response.FastCityRankingsResponse {
	response2 := &response.FastCityRankingsResponse{
		ElectionID:  electionPairID,
		Rankings:    make([]response.CityRankingItem, 0, len(rankings)),
		LastUpdated: time.Now(),
		TotalCities: uint64(len(rankings)),
	}

	for _, ranking := range rankings {
		response2.Rankings = append(response2.Rankings, response.CityRankingItem{
			CityName:          ranking.CityName,
			ConfirmedVotes:    ranking.ConfirmedVotes,
			UniqueVoters:      ranking.UniqueVoters,
			ParticipationRate: ranking.ParticipationRate,
			CityRank:          ranking.CityRank,
		})

		if response2.LastUpdated.Before(ranking.LastUpdated) {
			response2.LastUpdated = ranking.LastUpdated
		}
	}

	return response2
}

// Cache helper methods (run asynchronously)
func (m *Module) cacheElectionResults(ctx context.Context, key string, result *response.FastElectionResultsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache election results")
	}
}

func (m *Module) cacheCityResults(ctx context.Context, key string, result *response.FastCityResultsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache city results")
	}
}

func (m *Module) cacheElectionSummary(ctx context.Context, key string, result *response.FastElectionSummaryResponse) {
	if err := m.redisCache.Set(ctx, key, result, SummaryCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache election summary")
	}
}

func (m *Module) cacheCityRankings(ctx context.Context, key string, result *response.FastCityRankingsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache city rankings")
	}
}

func (m *Module) cacheAllElections(ctx context.Context, key string, result *response.FastAllElectionsSummaryResponse) {
	if err := m.redisCache.Set(ctx, key, result, AllElectionsCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache all elections")
	}
}

// Helper functions
func calculateCompletionRate(summary *model.ElectionSummary) float64 {
	if summary.TotalUniqueVoters > 0 {
		return math.Min(100.0, (float64(summary.TotalConfirmedVotes)/float64(summary.TotalUniqueVoters))*100)
	}
	return 0.0
}

func calculateParticipationRate(votes, voters uint64) float64 {
	if voters > 0 {
		return (float64(votes) / float64(voters)) * 100
	}
	return 0.0
}

func (m *Module) warmCache(ctx context.Context, activeElections []string) {
	for _, electionID := range activeElections {
		// Warm election results cache
		go func(id string) {
			if _, err := m.GetLiveElectionResultsWithCache(ctx, id); err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"election_id": id,
				}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm election cache")
			}
		}(electionID)

		// Warm election summary cache
		go func(id string) {
			if _, err := m.GetElectionSummaryWithCache(ctx, id); err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"election_id": id,
				}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm summary cache")
			}
		}(electionID)

		// Warm city rankings cache
		go func(id string) {
			if _, err := m.GetCityRankingsWithCache(ctx, id, 20); err != nil {
				log.WithFields(log.Fields{
					"error":       err,
					"election_id": id,
				}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm rankings cache")
			}
		}(electionID)
	}

	// Warm all elections summary
	go func() {
		if _, err := m.GetAllElectionsSummaryWithCache(ctx); err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to warm all elections cache")
		}
	}()
}
