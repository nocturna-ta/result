package live_result

import (
	"context"
	"github.com/nocturna-ta/golib/log"
	"github.com/nocturna-ta/result/internal/domain/model"
	"github.com/nocturna-ta/result/internal/usecases/response"
	"math"
	"time"
)

func (m *Module) buildElectionResultsResponse(liveResults []*model.LiveElectionResult, summary *model.ElectionSummary, cityRankings []*model.CityRanking, electionPairID string) *response.ElectionResultsResponse {
	response2 := &response.ElectionResultsResponse{
		ElectionID:    electionPairID,
		RegionResults: make([]response.RegionResultSummary, 0, len(liveResults)),
		TopCities:     make([]response.CityResultSummary, 0, len(cityRankings)),
	}

	// Calculate totals from live results for vote tracking
	var totalAttempts uint64
	for _, result := range liveResults {
		totalAttempts += result.TotalVotes
	}

	// Add enhanced region results with new percentage fields
	for _, result := range liveResults {
		response2.RegionResults = append(response2.RegionResults, response.RegionResultSummary{
			Region:                         result.Region,
			Votes:                          result.ConfirmedVotes,
			TotalAttempts:                  result.TotalVotes,                     // NEW
			Percentage:                     result.RegionalDistributionPercentage, // UPDATED: Was using wrong field
			VoteShareInRegion:              result.VoteShareInRegion,              // NEW: Who's winning in region
			RegionalDistributionPercentage: result.RegionalDistributionPercentage, // NEW: Where votes come from
			SuccessRate:                    result.SuccessPercentage,              // NEW: Processing success rate
		})
	}

	// Add enhanced summary data if available
	if summary != nil {
		response2.TotalVotes = summary.TotalConfirmedVotes
		response2.TotalVoters = summary.TotalUniqueVoters
		response2.TotalVoteAttempts = summary.TotalVoteAttempts
		response2.LastUpdated = summary.LastUpdated
		response2.OverallPercentage = summary.OverallSuccessRate
		response2.OverallVoteSharePercentage = summary.OverallVoteSharePercentage // NEW: True election percentage

		// Enhanced overall stats
		response2.OverallStats = response.ElectionStatsSummary{
			TotalVoters:                     summary.TotalUniqueVoters,
			TotalRegions:                    summary.TotalRegions,
			SuccessRate:                     summary.OverallSuccessRate,
			ActiveRegions:                   uint64(len(liveResults)),
			CompletionRate:                  calculateCompletionRate(summary),
			AverageSuccessRateAcrossRegions: calculateAverageSuccessRate(liveResults), // NEW
		}

		// Calculate votes per second
		if !summary.LastUpdated.IsZero() {
			timeSinceStart := time.Since(summary.LastUpdated.Add(-24 * time.Hour))
			if timeSinceStart.Seconds() > 0 {
				response2.OverallStats.VotesPerSecond = float64(summary.TotalConfirmedVotes) / timeSinceStart.Seconds()
			}
		}
	}

	// Add enhanced top cities with new distribution percentages
	for _, city := range cityRankings {
		// Calculate success rate from confirmed vs unique voters
		successRate := calculateCitySuccessRate(city.ConfirmedVotes, city.UniqueVoters)

		response2.TopCities = append(response2.TopCities, response.CityResultSummary{
			City:                       city.CityName,
			Votes:                      city.ConfirmedVotes,
			Voters:                     city.UniqueVoters,
			Attempts:                   city.UniqueVoters, // Assuming attempts = voters for ranking
			Percentage:                 city.ParticipationRate,
			VoteDistributionPercentage: city.VoteDistributionPercentage, // NEW
			Rank:                       city.CityRank,
			SuccessRate:                successRate, // NEW
		})
	}

	return response2
}

func (m *Module) buildCityResultsResponse(cityResults []*model.LiveCityResult, cityName string) *response.CityResultsResponse {
	response2 := &response.CityResultsResponse{
		CityName:        cityName,
		ElectionResults: make([]response.CityElectionResult, 0, len(cityResults)),
		LastUpdated:     time.Now(),
	}

	var totalVoters, totalVotes, totalAttempts uint64
	var leadingElection string
	var maxVotes uint64
	var totalSuccessRate float64

	for _, result := range cityResults {
		response2.ElectionResults = append(response2.ElectionResults, response.CityElectionResult{
			ElectionPairID:             result.ElectionPairID,
			Votes:                      result.ConfirmedVotes,
			Percentage:                 result.GetCandidatePercentageInCity(),
			VoteDistributionPercentage: calculateVoteDistributionForCity(result), // NEW: Calculate from model
			VoteSuccessRate:            result.GetVoteSuccessRate(),
		})

		totalVoters += result.TotalUniqueVoters
		totalVotes += result.ConfirmedVotes
		totalSuccessRate += result.GetVoteSuccessRate()

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
	response2.TotalAttempts = totalAttempts // NEW

	// Calculate enhanced city statistics
	avgDistributionPercentage := float64(100) / float64(len(cityResults)) // Equal distribution baseline
	if len(cityResults) > 0 {
		totalSuccessRate = totalSuccessRate / float64(len(cityResults))
	}

	response2.CityStats = response.CityStatistics{
		TotalElections:                uint64(len(cityResults)),
		LeadingElection:               leadingElection,
		ParticipationRate:             calculateParticipationRate(totalVotes, totalVoters),
		AverageVotesPerElection:       float64(totalVotes) / float64(len(cityResults)),
		TotalVoteAttempts:             totalAttempts,             // NEW
		OverallSuccessRate:            totalSuccessRate,          // NEW
		AverageDistributionPercentage: avgDistributionPercentage, // NEW
	}

	return response2
}

func (m *Module) buildCityRankingsResponse(rankings []*model.CityRanking, electionPairID string) *response.CityRankingsResponse {
	response2 := &response.CityRankingsResponse{
		ElectionID:  electionPairID,
		Rankings:    make([]response.CityRankingItem, 0, len(rankings)),
		LastUpdated: time.Now(),
		TotalCities: uint64(len(rankings)),
	}

	for _, ranking := range rankings {
		// Calculate total attempts (assuming some reasonable ratio)
		totalAttempts := calculateEstimatedAttempts(ranking.ConfirmedVotes, ranking.UniqueVoters)
		successRate := calculateCitySuccessRate(ranking.ConfirmedVotes, totalAttempts)

		response2.Rankings = append(response2.Rankings, response.CityRankingItem{
			CityName:                   ranking.CityName,
			ConfirmedVotes:             ranking.ConfirmedVotes,
			UniqueVoters:               ranking.UniqueVoters,
			TotalAttempts:              totalAttempts, // NEW
			ParticipationRate:          ranking.ParticipationRate,
			VoteDistributionPercentage: ranking.VoteDistributionPercentage, // NEW
			CityRank:                   ranking.CityRank,
			SuccessRate:                successRate, // NEW
		})

		if response2.LastUpdated.Before(ranking.LastUpdated) {
			response2.LastUpdated = ranking.LastUpdated
		}
	}

	return response2
}

// NEW: Enhanced function to build election summary response
func (m *Module) buildElectionSummaryResponse(summary *model.ElectionSummary) *response.ElectionSummaryResponse {
	return &response.ElectionSummaryResponse{
		ElectionID:                 summary.ElectionPairID,
		TotalUniqueVoters:          summary.TotalUniqueVoters,
		TotalRegions:               summary.TotalRegions,
		TotalConfirmedVotes:        summary.TotalConfirmedVotes,
		TotalVoteAttempts:          summary.TotalVoteAttempts, // NEW
		OverallSuccessRate:         summary.OverallSuccessRate,
		OverallVoteSharePercentage: summary.OverallVoteSharePercentage, // NEW
		LastUpdated:                summary.LastUpdated,
	}
}

// NEW: Enhanced function to build all elections summary with global stats
func (m *Module) buildAllElectionsSummaryResponse(summaries []*model.ElectionSummary) *response.AllElectionsSummaryResponse {
	elections := make([]response.ElectionSummaryResponse, 0, len(summaries))

	// Build individual election summaries
	for _, summary := range summaries {
		elections = append(elections, *m.buildElectionSummaryResponse(summary))
	}

	return &response.AllElectionsSummaryResponse{
		Elections:      elections,
		LastUpdated:    time.Now(),
		TotalElections: uint64(len(summaries)),
	}
}

func calculateAverageSuccessRate(liveResults []*model.LiveElectionResult) float64 {
	if len(liveResults) == 0 {
		return 0.0
	}

	var total float64
	for _, result := range liveResults {
		total += result.SuccessPercentage
	}
	return total / float64(len(liveResults))
}

func calculateCitySuccessRate(confirmedVotes, totalAttempts uint64) float64 {
	if totalAttempts == 0 {
		return 0.0
	}
	return (float64(confirmedVotes) / float64(totalAttempts)) * 100.0
}

func calculateVoteDistributionForCity(result *model.LiveCityResult) float64 {
	// This would typically come from the database, but we can calculate a placeholder
	// In a real scenario, you'd query for this candidate's total votes across all cities
	// For now, return a placeholder that could be populated by a separate query
	return 0.0 // TODO: Implement with actual distribution calculation
}

func calculateEstimatedAttempts(confirmedVotes, uniqueVoters uint64) uint64 {
	// Estimate total attempts based on confirmed votes and unique voters
	// Assuming some users might vote multiple times or have failed attempts
	if uniqueVoters > confirmedVotes {
		return uniqueVoters + (uniqueVoters-confirmedVotes)/10 // Add 10% estimate for failed attempts
	}
	return confirmedVotes + confirmedVotes/20 // Add 5% estimate for failed attempts
}

// Cache helper methods (run asynchronously)
func (m *Module) cacheElectionResults(ctx context.Context, key string, result *response.ElectionResultsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache election results")
	}
}

func (m *Module) cacheCityResults(ctx context.Context, key string, result *response.CityResultsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache city results")
	}
}

func (m *Module) cacheElectionSummary(ctx context.Context, key string, result *response.ElectionSummaryResponse) {
	if err := m.redisCache.Set(ctx, key, result, SummaryCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache election summary")
	}
}

func (m *Module) cacheCityRankings(ctx context.Context, key string, result *response.CityRankingsResponse) {
	if err := m.redisCache.Set(ctx, key, result, DefaultCacheTTL); err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"key":   key,
		}).WarnWithCtx(ctx, "[FastLiveResultUseCase] Failed to cache city rankings")
	}
}

func (m *Module) cacheAllElections(ctx context.Context, key string, result *response.AllElectionsSummaryResponse) {
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
