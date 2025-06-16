package response

import "time"

type ElectionLiveResult struct {
	ElectionPairID string    `json:"election_pair_id"`
	TotalVotes     uint64    `json:"total_votes"`
	Percentage     float64   `json:"percentage"`
	LastUpdated    time.Time `json:"last_updated"`
}

type ElectionOverallResult struct {
	ElectionPairID string                `json:"election_pair_id"`
	TotalVotes     uint64                `json:"total_votes"`
	TotalVoters    uint64                `json:"total_voters"`
	Candidates     []*ElectionLiveResult `json:"candidates"`
	LastUpdated    time.Time             `json:"last_updated"`
}

type CityLiveResult struct {
	CityName    string             `json:"city_name"`
	TotalVoters uint64             `json:"total_voters"`
	TotalVotes  uint64             `json:"total_votes"`
	Candidates  []*CandidateResult `json:"candidates"`
	LastUpdated time.Time          `json:"last_updated"`
}

type CandidateResult struct {
	ElectionPairID string  `json:"election_pair_id"`
	VoteCount      uint64  `json:"vote_count"`
	Percentage     float64 `json:"percentage"`
}

type LiveResultsSummary struct {
	OverallResults []*ElectionOverallResult `json:"overall_results"`
	CityResults    []*CityLiveResult        `json:"city_results"`
	LastUpdated    time.Time                `json:"last_updated"`
	TotalProcessed uint64                   `json:"total_processed"`
}
