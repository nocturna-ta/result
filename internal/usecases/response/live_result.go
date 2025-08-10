package response

import "time"

type ElectionResultsResponse struct {
	ElectionID                 string                `json:"election_id"`
	TotalVotes                 uint64                `json:"total_votes"`
	TotalVoters                uint64                `json:"total_voters"`
	TotalVoteAttempts          uint64                `json:"total_vote_attempts"` // NEW
	LastUpdated                time.Time             `json:"last_updated"`
	OverallPercentage          float64               `json:"overall_percentage"`
	OverallVoteSharePercentage float64               `json:"overall_vote_share_percentage"` // NEW
	RegionResults              []RegionResultSummary `json:"regions"`
	TopCities                  []CityResultSummary   `json:"top_cities,omitempty"`
	OverallStats               ElectionStatsSummary  `json:"stats"`
}

type CityResultsResponse struct {
	CityName        string               `json:"city_name"`
	TotalVoters     uint64               `json:"total_voters"`
	TotalVotes      uint64               `json:"total_votes"`
	TotalAttempts   uint64               `json:"total_attempts"` // NEW
	LastUpdated     time.Time            `json:"last_updated"`
	ElectionResults []CityElectionResult `json:"election_results"`
	CityStats       CityStatistics       `json:"city_stats"`
}

type ElectionSummaryResponse struct {
	ElectionID                 string    `json:"election_id"`
	TotalUniqueVoters          uint64    `json:"total_unique_voters"`
	TotalRegions               uint64    `json:"total_regions"`
	TotalConfirmedVotes        uint64    `json:"total_confirmed_votes"`
	TotalVoteAttempts          uint64    `json:"total_vote_attempts"` // NEW
	OverallSuccessRate         float64   `json:"overall_success_rate"`
	OverallVoteSharePercentage float64   `json:"overall_vote_share_percentage"` // NEW
	LastUpdated                time.Time `json:"last_updated"`
}

type CityRankingsResponse struct {
	ElectionID  string            `json:"election_id"`
	Rankings    []CityRankingItem `json:"rankings"`
	LastUpdated time.Time         `json:"last_updated"`
	TotalCities uint64            `json:"total_cities"`
}

type AllElectionsSummaryResponse struct {
	Elections      []ElectionSummaryResponse `json:"elections"`
	LastUpdated    time.Time                 `json:"last_updated"`
	TotalElections uint64                    `json:"total_elections"`
}

// Enhanced sub-structs with new fields

type RegionResultSummary struct {
	Region                         string  `json:"region"`
	Votes                          uint64  `json:"votes"`
	TotalAttempts                  uint64  `json:"total_attempts"` // NEW
	Percentage                     float64 `json:"percentage"`
	VoteShareInRegion              float64 `json:"vote_share_in_region"`             // NEW
	RegionalDistributionPercentage float64 `json:"regional_distribution_percentage"` // NEW
	SuccessRate                    float64 `json:"success_rate"`                     // NEW
}

type CityResultSummary struct {
	City                       string  `json:"city"`
	Votes                      uint64  `json:"votes"`
	Voters                     uint64  `json:"voters"`
	Attempts                   uint64  `json:"attempts"` // NEW
	Percentage                 float64 `json:"percentage"`
	VoteDistributionPercentage float64 `json:"vote_distribution_percentage"` // NEW
	Rank                       uint32  `json:"rank"`
	SuccessRate                float64 `json:"success_rate"` // NEW
}

type ElectionStatsSummary struct {
	TotalVoters                     uint64  `json:"total_voters"`
	TotalRegions                    uint64  `json:"total_regions"`
	SuccessRate                     float64 `json:"success_rate"`
	VotesPerSecond                  float64 `json:"votes_per_second"`
	ActiveRegions                   uint64  `json:"active_regions"`
	CompletionRate                  float64 `json:"completion_rate"`
	AverageSuccessRateAcrossRegions float64 `json:"avg_success_rate_across_regions"` // NEW
}

type CityElectionResult struct {
	ElectionPairID             string  `json:"election_pair_id"`
	Votes                      uint64  `json:"votes"`
	Attempts                   uint64  `json:"attempts"` // NEW
	Percentage                 float64 `json:"percentage"`
	VoteDistributionPercentage float64 `json:"vote_distribution_percentage"` // NEW
	CityRank                   uint32  `json:"city_rank"`
	VoteSuccessRate            float64 `json:"vote_success_rate"`
}

type CityStatistics struct {
	ParticipationRate             float64 `json:"participation_rate"`
	AverageVotesPerElection       float64 `json:"avg_votes_per_election"`
	TotalElections                uint64  `json:"total_elections"`
	LeadingElection               string  `json:"leading_election"`
	TotalVoteAttempts             uint64  `json:"total_vote_attempts"`         // NEW
	OverallSuccessRate            float64 `json:"overall_success_rate"`        // NEW
	AverageDistributionPercentage float64 `json:"avg_distribution_percentage"` // NEW
}

type CityRankingItem struct {
	CityName                   string  `json:"city_name"`
	ConfirmedVotes             uint64  `json:"confirmed_votes"`
	UniqueVoters               uint64  `json:"unique_voters"`
	TotalAttempts              uint64  `json:"total_attempts"` // NEW
	ParticipationRate          float64 `json:"participation_rate"`
	VoteDistributionPercentage float64 `json:"vote_distribution_percentage"` // NEW
	CityRank                   uint32  `json:"city_rank"`
	SuccessRate                float64 `json:"success_rate"` // NEW
}

type IncrementalUpdateData struct {
	Type       string                 `json:"type"`
	ElectionID string                 `json:"election_id"`
	Region     string                 `json:"region,omitempty"`
	CityName   string                 `json:"city_name,omitempty"`
	NewVotes   uint64                 `json:"new_votes"`
	Changes    map[string]interface{} `json:"changes"`
	Timestamp  time.Time              `json:"timestamp"`
}

type CacheStatisticsResponse struct {
	HitRate          float64                `json:"hit_rate"`
	TotalHits        uint64                 `json:"total_hits"`
	TotalMisses      uint64                 `json:"total_misses"`
	CacheSize        string                 `json:"cache_size"`
	ConnectedClients int                    `json:"connected_clients"`
	LastCacheRefresh time.Time              `json:"last_cache_refresh"`
	ActiveElections  []string               `json:"active_elections"`
	RedisInfo        map[string]interface{} `json:"redis_info"`
}

type VoteResultResponse struct {
	ID              string     `json:"id"`
	VoterID         string     `json:"voter_id"`
	ElectionPairID  string     `json:"election_pair_id"`
	Region          string     `json:"region"`
	Status          string     `json:"status"`
	TransactionHash string     `json:"transaction_hash"`
	ErrorMessage    string     `json:"error_message,omitempty"`
	VotedAt         time.Time  `json:"voted_at"`
	ProcessedAt     *time.Time `json:"processed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
